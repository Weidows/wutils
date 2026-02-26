package zip

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/yeka/zip"
)

// Result represents the result of an unzip operation
type Result struct {
	Success   bool
	Extracted []string // List of extracted file paths
	Error     error
}

// Archive represents a compressed archive with optional password protection
type Archive struct {
	archivePath string
	password    string
}

// NewArchive creates a new Archive instance
func NewArchive(archivePath, password string) *Archive {
	return &Archive{
		archivePath: archivePath,
		password:    password,
	}
}

// Unzip extracts all files from the archive to the specified destination directory
// Returns a Result with extracted file paths and any error encountered
func (a *Archive) Unzip(destDir string) Result {
	if destDir == "" {
		destDir = "."
	}

	ext := strings.ToLower(filepath.Ext(a.archivePath))
	switch ext {
	case ".zip":
		return a.unzipZip(destDir)
	case ".7z":
		return a.unzip7z(destDir)
	default:
		return Result{
			Success: false,
			Error:   fmt.Errorf("unsupported archive format: %s", ext),
		}
	}
}

// TryUnzip attempts to unzip the archive with the given password
// Returns true if successful, false otherwise (deprecated: use Unzip instead)
func (a *Archive) TryUnzip() bool {
	ext := strings.ToLower(filepath.Ext(a.archivePath))
	switch ext {
	case ".zip":
		reader, err := zip.OpenReader(a.archivePath)
		if err != nil {
			return false
		}
		defer reader.Close()

		for _, file := range reader.File {
			if err := a.extractZipFile(file, "."); err != nil {
				return false
			}
			return true
		}

	case ".7z":
		reader, err := sevenzip.OpenReaderWithPassword(a.archivePath, a.password)
		if err != nil {
			return false
		}
		defer reader.Close()

		for _, file := range reader.File {
			if err := a.extract7zFile(file, "."); err != nil {
				return false
			}
			return true
		}
	}

	return false
}

// unzipZip handles ZIP archive extraction
func (a *Archive) unzipZip(destDir string) Result {
	reader, err := zip.OpenReader(a.archivePath)
	if err != nil {
		return Result{Success: false, Error: fmt.Errorf("failed to open zip: %w", err)}
	}
	defer reader.Close()

	return a.extractAllZip(reader, destDir)
}

// unzip7z handles 7z archive extraction
func (a *Archive) unzip7z(destDir string) Result {
	reader, err := sevenzip.OpenReaderWithPassword(a.archivePath, a.password)
	if err != nil {
		return Result{Success: false, Error: fmt.Errorf("failed to open 7z: %w", err)}
	}
	defer reader.Close()

	return a.extractAll7z(reader, destDir)
}

// extractAllZip extracts all files from a ZIP archive
func (a *Archive) extractAllZip(reader *zip.ReadCloser, destDir string) Result {
	var extracted []string

	for _, file := range reader.File {
		if err := a.extractZipFile(file, destDir); err != nil {
			return Result{
				Success:   false,
				Extracted: extracted,
				Error:     fmt.Errorf("failed to extract %s: %w", file.Name, err),
			}
		}
		extracted = append(extracted, filepath.Join(destDir, file.Name))
	}

	return Result{Success: true, Extracted: extracted}
}

// extractAll7z extracts all files from a 7z archive
func (a *Archive) extractAll7z(reader *sevenzip.ReadCloser, destDir string) Result {
	var extracted []string

	for _, file := range reader.File {
		if err := a.extract7zFile(file, destDir); err != nil {
			return Result{
				Success:   false,
				Extracted: extracted,
				Error:     fmt.Errorf("failed to extract %s: %w", file.Name, err),
			}
		}
		extracted = append(extracted, filepath.Join(destDir, file.Name))
	}

	return Result{Success: true, Extracted: extracted}
}

// extractZipFile extracts a single file from a ZIP archive
func (a *Archive) extractZipFile(file *zip.File, destDir string) error {
	if file.FileInfo().IsDir() {
		return nil
	}

	if file.IsEncrypted() {
		file.SetPassword(a.password)
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in archive: %w", err)
	}
	defer src.Close()

	destPath := filepath.Join(destDir, file.Name)
	dirPath := filepath.Dir(destPath)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	return err
}

// extract7zFile extracts a single file from a 7z archive
func (a *Archive) extract7zFile(file *sevenzip.File, destDir string) error {
	fi := file.FileInfo()
	if fi.IsDir() {
		return nil
	}

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file in archive: %w", err)
	}
	defer src.Close()

	destPath := filepath.Join(destDir, file.Name)
	dirPath := filepath.Dir(destPath)

	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	dest, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dest.Close()

	_, err = io.Copy(dest, src)
	return err
}

// VerifyPassword checks if the password is correct without extracting files
func (a *Archive) VerifyPassword() bool {
	ext := strings.ToLower(filepath.Ext(a.archivePath))
	switch ext {
	case ".zip":
		reader, err := zip.OpenReader(a.archivePath)
		if err != nil {
			return false
		}
		defer reader.Close()

		if len(reader.File) == 0 {
			return a.password == ""
		}

		for _, file := range reader.File {
			file.SetPassword(a.password)
			rc, err := file.Open()
			if err != nil {
				return false
			}
			buf := make([]byte, 1)
			_, err = rc.Read(buf)
			rc.Close()
			if err != nil {
				return false
			}
		}
		return true

	case ".7z":
		reader, err := sevenzip.OpenReaderWithPassword(a.archivePath, a.password)
		if err != nil {
			return false
		}
		defer reader.Close()
		return true
	}

	return false
}
