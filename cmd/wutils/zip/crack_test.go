package zip

import (
	"testing"
)

func TestCrackPassword(t *testing.T) {
	type args struct {
		archivePath string
		passwords   []string
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
				passwords:   []string{"test", "wrong", "123456"},
			},
			want: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CrackPasswordWithList(tt.args.archivePath, tt.args.passwords); got != tt.want {
				t.Errorf("CrackPassword() = %v, want %v", got, tt.want)
			}
		})
	}
}
