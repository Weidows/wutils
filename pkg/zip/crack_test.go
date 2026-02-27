package zip

import (
	"os"
	"path/filepath"
	"testing"
)

func TestArchive_TryUnzip(t *testing.T) {
	tests := []struct {
		name    string
		archive *Archive
		want    bool
	}{
		{
			name:    "test-zip-right-password",
			archive: &Archive{archivePath: "test.zip", password: "test"},
			want:    true,
		},
		{
			name:    "test-7z-right-password",
			archive: &Archive{archivePath: "test.7z", password: "test"},
			want:    true,
		},
		{
			name:    "test-nonexistent-file",
			archive: &Archive{archivePath: "nonexistent.zip", password: ""},
			want:    false,
		},
		{
			name:    "test-unsupported-format",
			archive: &Archive{archivePath: "test.rar", password: ""},
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.archive.TryUnzip(); got != tt.want {
				t.Errorf("Archive.TryUnzip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchive_VerifyPassword(t *testing.T) {
	tests := []struct {
		name    string
		archive *Archive
		want    bool
	}{
		{
			name:    "zip-valid-password",
			archive: NewArchive("test.zip", "test"),
			want:    true,
		},
		{
			name:    "7z-valid-password",
			archive: NewArchive("test.7z", "test"),
			want:    true,
		},
		{
			name:    "nonexistent-file",
			archive: NewArchive("nonexistent.zip", ""),
			want:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.archive.VerifyPassword(); got != tt.want {
				t.Errorf("Archive.VerifyPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchive_Unzip(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		archive *Archive
		wantErr bool
	}{
		{
			name:    "unzip-valid-zip",
			archive: NewArchive("test.zip", "test"),
			wantErr: false,
		},
		{
			name:    "unzip-valid-7z",
			archive: NewArchive("test.7z", "test"),
			wantErr: false,
		},
		{
			name:    "unzip-wrong-password",
			archive: NewArchive("test.zip", "wrong"),
			wantErr: true,
		},
		{
			name:    "unzip-nonexistent",
			archive: NewArchive("nonexistent.zip", ""),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.archive.Unzip(tmpDir)
			if (result.Error != nil) != tt.wantErr {
				t.Errorf("Archive.Unzip() error = %v, wantErr %v", result.Error, tt.wantErr)
			}
		})
	}
}

func TestNewArchive(t *testing.T) {
	archive := NewArchive("test.zip", "password")

	if archive.archivePath != "test.zip" {
		t.Errorf("expected archivePath to be 'test.zip', got %s", archive.archivePath)
	}
	if archive.password != "password" {
		t.Errorf("expected password to be 'password', got %s", archive.password)
	}
}

func TestArchive_ExtractToCustomDir(t *testing.T) {
	tmpDir := t.TempDir()
	customDir := filepath.Join(tmpDir, "output")

	archive := NewArchive("test.zip", "test")
	result := archive.Unzip(customDir)

	if result.Error != nil {
		t.Errorf("Unzip() failed: %v", result.Error)
	}

	if !result.Success {
		t.Error("Expected success=true")
	}

	if len(result.Extracted) == 0 {
		t.Error("Expected at least one extracted file")
	}

	for _, path := range result.Extracted {
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected file %s to exist", path)
		}
	}
}
