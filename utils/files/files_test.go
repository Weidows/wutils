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

func TestIsExist(t *testing.T) {
	fmt.Println(IsExist("./2e3ce48952af857ccbecb2e8f7ff52c6.mp4"))
}

func TestCopyFile(t *testing.T) {
	type args struct {
		src string
		dst string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"right",
			args{
				"2e3ce48952af857ccbecb2e8f7ff52c6.mp4",
				"2e3ce48952af857ccbecb2e8f7ff52c6_cp.mp4",
			},
			false,
		},
		{
			"wrong",
			args{
				"2e3ce48952af857ccbecb2e8f7ff52c6_.mp4",
				"2e3ce48952af857ccbecb2e8f7ff52c6_cp.mp4",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CopyFile(tt.args.src, tt.args.dst); (err != nil) != tt.wantErr {
				t.Errorf("CopyFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
