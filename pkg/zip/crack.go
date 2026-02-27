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

// Archive represents a password-protected archive
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

// TryUnzip verifies if the password is correct without extracting files
// Returns true if successful, false otherwise
func (a *Archive) TryUnzip() bool {
	ext := strings.ToLower(filepath.Ext(a.archivePath))

	switch ext {
	case ".zip":
		return a.tryUnzipZip()
	case ".7z":
		return a.tryUnzip7z()
	default:
		return false
	}
}

func (a *Archive) tryUnzipZip() bool {
	r, err := zip.OpenReader(a.archivePath)
	if err != nil {
		return false
	}
	defer r.Close()

	// Check if the zip file is password protected
	for _, f := range r.File {
		if f.IsEncrypted() {
			// Set the password and try to read
			f.SetPassword(a.password)
			rc, err := f.Open()
			if err != nil {
				// Password is wrong
				return false
			}
			// Read content to verify password - wrong password produces garbage
			content, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				// Password is wrong or error reading
				return false
			}
			// Check if content is garbage (wrong password)
			if len(content) > 0 && a.isGarbageContent(content) {
				return false
			}
			return true
		}
	}

	// If not encrypted, just check if we can read the file
	return true
}

func (a *Archive) tryUnzip7z() bool {
	// Check if file exists first
	if _, err := os.Stat(a.archivePath); os.IsNotExist(err) {
		return false
	}

	r, err := sevenzip.OpenReader(a.archivePath)
	if err != nil {
		return false
	}
	defer r.Close()

	// Try to read first file content to verify password
	// The bodgit/sevenzip library doesn't support password explicitly,
	// but will fail when trying to read encrypted content without proper key
	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			// Password required or wrong
			return false
		}
		// Try to read some content
		buf := make([]byte, 32)
		n, err := rc.Read(buf)
		rc.Close()
		if err != nil && n == 0 {
			// Can't read - likely encrypted/wrong password
			return false
		}
		// If we can read, password is correct
		return true
	}

	return true
}

// Unzip extracts the archive to the specified directory
func (a *Archive) Unzip(dest string) error {
	ext := strings.ToLower(filepath.Ext(a.archivePath))

	switch ext {
	case ".zip":
		return a.unzipZip(dest)
	case ".7z":
		return a.unzip7z(dest)
	default:
		return fmt.Errorf("unsupported archive format: %s", ext)
	}
}

func (a *Archive) unzipZip(dest string) error {
	r, err := zip.OpenReader(a.archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		if f.IsEncrypted() {
			f.SetPassword(a.password)
		}

		path := filepath.Join(dest, f.Name)

		// Check for zip slip vulnerability
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to open file %s: %w (wrong password?)", f.Name, err)
		}

		// Read all content to verify password is correct
		// Wrong password may extract garbage without returning error
		content, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			outFile.Close()
			return fmt.Errorf("failed to read file %s: %w (wrong password?)", f.Name, err)
		}

		// Check if content is likely garbage (wrong password produces random bytes)
		// A properly decrypted file should have some structure or not be all random
		if len(content) > 0 && a.isGarbageContent(content) {
			outFile.Close()
			return fmt.Errorf("file %s content is corrupted (wrong password)", f.Name)
		}

		_, err = outFile.Write(content)
		outFile.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// isGarbageContent checks if the extracted content is likely garbage
// Wrong password decryption produces random-looking bytes
func (a *Archive) isGarbageContent(content []byte) bool {
	// If no password was provided and file is encrypted, content is likely garbage
	if a.password == "" {
		return true
	}

	// Check for high entropy (random bytes) - typical of wrong password
	// Count byte value distribution
	counts := [256]int{}
	for _, b := range content {
		counts[b]++
	}

	// For wrong password, byte distribution should be relatively uniform
	// Count how many byte values appear
	uniqueBytes := 0
	for _, c := range counts {
		if c > 0 {
			uniqueBytes++
		}
	}

	// If more than 200 different byte values in a small sample, it's likely random
	if len(content) >= 256 && uniqueBytes > 200 {
		return true
	}

	// Check for null bytes (common in garbage)
	nullRatio := float64(counts[0]) / float64(len(content))
	if nullRatio > 0.5 && len(content) > 100 {
		return true
	}

	return false
}

func (a *Archive) unzip7z(dest string) error {
	r, err := sevenzip.OpenReader(a.archivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)

		// Check for zip slip vulnerability
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		_, err = io.Copy(outFile, rc)
		if err != nil {
			return err
		}
	}

	return nil
}
