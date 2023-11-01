package files

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestGetSubFiles(t *testing.T) {
	files := GetAllSubFiles("../../utils")
	for _, v := range files {
		fmt.Println(filepath.Join(v.Path, v.Name))
	}
}

func TestGetSubFilesWithFilter(t *testing.T) {
	files := GetAllSubFilesWithFilter("../../utils", func(fileInfo *subFileInfo) bool {
		if filepath.Ext(fileInfo.Name) == ".go" {
			return true
		}
		return false
	})
	for _, v := range files {
		fmt.Println(filepath.Join(v.Path, v.Name))
	}
}
