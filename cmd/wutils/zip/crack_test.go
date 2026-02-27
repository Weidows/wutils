package zip

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCrackPassword(t *testing.T) {
	localDictPath := "./password-dict.txt"
	os.WriteFile(localDictPath, []byte("test\nwrong\n123456"), 0644)
	defer os.Remove(localDictPath)

	home, _ := os.UserHomeDir()
	homeDictPath := filepath.Join(home, ".config", "wutils", "password-dict.txt")
	os.WriteFile(homeDictPath, []byte("test\nwrong\n123456"), 0644)

	type args struct {
		archivePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test",
			args: args{
				archivePath: "test/test.zip",
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CrackPassword(tt.args.archivePath); got != tt.want {
				t.Errorf("CrackPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
