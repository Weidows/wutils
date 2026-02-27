package zip

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCrackPassword(t *testing.T) {
	home, _ := os.UserHomeDir()
	configDir := filepath.Join(home, ".config", "wutils")
	passwordDictPath := filepath.Join(configDir, "password-dict.txt")

	os.MkdirAll(configDir, 0755)
	os.WriteFile(passwordDictPath, []byte("test\nwrong\n123456"), 0644)

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
				archivePath: "test.zip",
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
