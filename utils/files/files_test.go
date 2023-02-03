package files

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestGetSubFiles(t *testing.T) {
	files := GetSubFiles("C:\\Users\\Administrator\\.config")
	for _, v := range files {
		fmt.Println(filepath.Join(v.Path, v.Name))
	}
}

func TestGetSubFilesWithFilter(t *testing.T) {
	files := GetSubFilesWithFilter("C:\\Users\\Administrator\\.config", func(fileInfo *subFileInfo) bool {
		if filepath.Ext(fileInfo.Name) == ".log" {
			return true
		}
		return false
	})
	for _, v := range files {
		fmt.Println(filepath.Join(v.Path, v.Name))
	}
}
