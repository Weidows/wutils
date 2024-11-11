package zip

import (
	"bytes"
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
	buffer := new(bytes.Buffer)

	// extracted, err := os.Create(fileName)
	// if err != nil {
	// 	return false
	// }
	// defer extracted.Close()

	// // 无论是否解压成功, 都要删除解压文件
	// defer os.Remove(fileName)

	n, err := io.Copy(buffer, zipped)
	if err != nil || n == 0 {
		return false
	}

	return true
}
