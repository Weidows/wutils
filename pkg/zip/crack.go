package zip

import (
	"archive/tar"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/nwaples/rardecode/v2"
	"github.com/yeka/zip"
)

type archiveType string

const (
	archiveTypeUnknown archiveType = "unknown"
	archiveTypeZip     archiveType = "zip"
	archiveType7z      archiveType = "7z"
	archiveTypeRAR     archiveType = "rar"
	archiveTypeTarGz   archiveType = "tar.gz"
	archiveTypeTarBz2  archiveType = "tar.bz2"
)

var (
	reSplit7z = regexp.MustCompile(`(?i)^(.*\.7z)\.(\d{3})$`)
	reZipVol  = regexp.MustCompile(`(?i)^\.z\d{2}$`)
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
	archivePath := ResolveArchivePath(a.archivePath)
	typeName := detectArchiveType(archivePath)

	switch typeName {
	case archiveTypeZip:
		return a.tryUnzipZip(archivePath)
	case archiveType7z:
		return a.tryUnzip7z(archivePath)
	case archiveTypeRAR:
		return a.tryUnzipRAR(archivePath)
	case archiveTypeTarGz:
		return a.tryUnzipTarGz(archivePath)
	case archiveTypeTarBz2:
		return a.tryUnzipTarBz2(archivePath)
	default:
		return false
	}
}

func (a *Archive) tryUnzipZip(archivePath string) bool {
	r, err := zip.OpenReader(archivePath)
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

func (a *Archive) tryUnzip7z(archivePath string) bool {
	// Check if file exists first
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		return false
	}

	r, err := sevenzip.OpenReader(archivePath)
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

func (a *Archive) tryUnzipRAR(archivePath string) bool {
	r, err := rardecode.OpenReader(archivePath, rardecode.Password(a.password))
	if err != nil {
		return false
	}
	defer r.Close()

	for {
		_, err := r.Next()
		if err == io.EOF {
			return true
		}
		if err != nil {
			return false
		}

		buf := make([]byte, 32)
		n, readErr := r.Read(buf)
		if readErr != nil && readErr != io.EOF {
			return false
		}
		if n > 0 || readErr == io.EOF {
			return true
		}
	}
}

func (a *Archive) tryUnzipTarGz(archivePath string) bool {
	f, err := os.Open(archivePath)
	if err != nil {
		return false
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return false
	}
	defer gzr.Close()

	return a.tryReadTarStream(gzr)
}

func (a *Archive) tryUnzipTarBz2(archivePath string) bool {
	f, err := os.Open(archivePath)
	if err != nil {
		return false
	}
	defer f.Close()

	bzr := bzip2.NewReader(f)
	return a.tryReadTarStream(bzr)
}

func (a *Archive) tryReadTarStream(r io.Reader) bool {
	tr := tar.NewReader(r)
	for {
		_, err := tr.Next()
		if err == io.EOF {
			return true
		}
		if err != nil {
			return false
		}

		buf := make([]byte, 32)
		n, readErr := tr.Read(buf)
		if readErr != nil && readErr != io.EOF {
			return false
		}
		if n > 0 || readErr == io.EOF {
			return true
		}
	}
}

// Unzip extracts the archive to the specified directory
func (a *Archive) Unzip(dest string) error {
	archivePath := ResolveArchivePath(a.archivePath)
	typeName := detectArchiveType(archivePath)

	switch typeName {
	case archiveTypeZip:
		return a.unzipZip(archivePath, dest)
	case archiveType7z:
		return a.unzip7z(archivePath, dest)
	case archiveTypeRAR:
		return a.unzipRAR(archivePath, dest)
	case archiveTypeTarGz:
		return a.unzipTarGz(archivePath, dest)
	case archiveTypeTarBz2:
		return a.unzipTarBz2(archivePath, dest)
	default:
		return fmt.Errorf("unsupported archive format: %s", archivePath)
	}
}

func (a *Archive) unzipZip(archivePath, dest string) error {
	r, err := zip.OpenReader(archivePath)
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

func (a *Archive) unzip7z(archivePath, dest string) error {
	r, err := sevenzip.OpenReader(archivePath)
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

func (a *Archive) unzipRAR(archivePath, dest string) error {
	r, err := rardecode.OpenReader(archivePath, rardecode.Password(a.password))
	if err != nil {
		return err
	}
	defer r.Close()

	for {
		hdr, err := r.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dest, hdr.Name)
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", hdr.Name)
		}

		if hdr.IsDir {
			if err := os.MkdirAll(path, hdr.Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, hdr.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, r)
		closeErr := outFile.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
	}
}

func (a *Archive) unzipTarGz(archivePath, dest string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	gzr, err := gzip.NewReader(f)
	if err != nil {
		return err
	}
	defer gzr.Close()

	return unzipTarStream(gzr, dest)
}

func (a *Archive) unzipTarBz2(archivePath, dest string) error {
	f, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer f.Close()

	bzr := bzip2.NewReader(f)
	return unzipTarStream(bzr, dest)
}

func unzipTarStream(r io.Reader, dest string) error {
	tr := tar.NewReader(r)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dest, hdr.Name)
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", hdr.Name)
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, os.FileMode(hdr.Mode)); err != nil {
				return err
			}
		case tar.TypeReg, tar.TypeRegA:
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}

			outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, tr)
			closeErr := outFile.Close()
			if err != nil {
				return err
			}
			if closeErr != nil {
				return closeErr
			}
		}
	}
}

func ResolveArchivePath(archivePath string) string {
	lowerPath := strings.ToLower(archivePath)

	if matches := reSplit7z.FindStringSubmatch(lowerPath); len(matches) == 3 {
		baseWith7z := archivePath[:len(matches[1])]
		return baseWith7z + ".001"
	}

	ext := strings.ToLower(filepath.Ext(archivePath))
	if reZipVol.MatchString(ext) {
		return strings.TrimSuffix(archivePath, filepath.Ext(archivePath)) + ".zip"
	}

	return archivePath
}

func detectArchiveType(archivePath string) archiveType {
	lowerPath := strings.ToLower(archivePath)

	switch {
	case strings.HasSuffix(lowerPath, ".tar.gz"):
		return archiveTypeTarGz
	case strings.HasSuffix(lowerPath, ".tar.bz2"):
		return archiveTypeTarBz2
	case strings.HasSuffix(lowerPath, ".zip"):
		return archiveTypeZip
	case strings.HasSuffix(lowerPath, ".7z") || reSplit7z.MatchString(lowerPath):
		return archiveType7z
	case strings.HasSuffix(lowerPath, ".rar"):
		return archiveTypeRAR
	default:
		return archiveTypeUnknown
	}
}
