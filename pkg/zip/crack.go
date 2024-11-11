package zip

import (
	"io"
	"path/filepath"
	"strings"

	"github.com/bodgit/sevenzip"
	"github.com/yeka/zip"
)

type Archive struct {
	archivePath string
	password    string
}

func NewArchive(archivePath, password string) *Archive {
	return &Archive{archivePath: archivePath, password: password}
}

// Attempt to unlock the archive with the given password
func (a *Archive) TryUnzip() bool {
	switch strings.ToLower(filepath.Ext(a.archivePath)) {
	case ".zip":
		reader, err := zip.OpenReader(a.archivePath)
		if err != nil {
			return false
		}
		defer reader.Close()

		for _, file := range reader.File {
			file.SetPassword(a.password)
			zipped, err := file.Open()
			if err != nil {
				return false
			}
			return unzip(zipped, file.Name)
		}

	case ".7z":
		reader, err := sevenzip.OpenReaderWithPassword(a.archivePath, a.password)
		if err != nil {
			return false
		}
		defer reader.Close()
		for _, file := range reader.File {
			zipped, err := file.Open()
			if err != nil {
				return false
			}
			return unzip(zipped, file.Name)
		}

	default:
		return false
	}

	return true
}

func unzip(zipped io.ReadCloser, fileName string) bool {
	defer zipped.Close()
	// 使用流式处理，避免将整个文件内容加载到内存中
	buffer := make([]byte, 1024) // 使用较小的缓冲区
	n, err := zipped.Read(buffer)
	if (err != nil && err != io.EOF) || n == 0 {
		return false
	}

	return true
}
