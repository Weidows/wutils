package zip

import (
	"testing"
)

func TestArchive_TryUnzip(t *testing.T) {
	tests := []struct {
		name    string
		archive *Archive
		want    bool
	}{
		{
			name:    "zip-wrong-password",
			archive: NewArchive("test.zip", "wrong"),
			want:    false,
		},
		{
			name:    "7z-wrong-password",
			archive: NewArchive("test.7z", "wrong"),
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

func TestArchive_Unzip(t *testing.T) {
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
		// Skip 7z tests as test file doesn't exist
		// {
		// 	name:    "unzip-valid-7z",
		// 	archive: NewArchive("test.7z", "test"),
		// 	wantErr: false,
		// },
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
			err := tt.archive.Unzip("./testoutput")
			if (err != nil) != tt.wantErr {
				t.Errorf("Archive.Unzip() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewArchive(t *testing.T) {
	archive := NewArchive("test.zip", "password")
	if archive == nil {
		t.Error("NewArchive() returned nil")
	}
	if archive.archivePath != "test.zip" {
		t.Errorf("NewArchive().archivePath = %v, want test.zip", archive.archivePath)
	}
	if archive.password != "password" {
		t.Errorf("NewArchive().password = %v, want password", archive.password)
	}
}
