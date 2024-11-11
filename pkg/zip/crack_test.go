package zip

import "testing"

func TestArchive_TryUnlock(t *testing.T) {
	tests := []struct {
		name string
		a    *Archive
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test-zip-wrong",
			a:    &Archive{archivePath: "test.zip", password: "test_wrong"},
			want: false,
		},
		{
			name: "test-zip-right",
			a:    &Archive{archivePath: "test.zip", password: "test"},
			want: true,
		},
		{
			name: "test-7z-wrong",
			a:    &Archive{archivePath: "test.7z", password: "test_wrong"},
			want: false,
		},
		{
			name: "test-7z-right",
			a:    &Archive{archivePath: "test.7z", password: "test"},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.a.TryUnzip(); got != tt.want {
				t.Errorf("Archive.tryUnlock() = %v, want %v", got, tt.want)
			}
		})
	}
}
