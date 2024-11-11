package zip

import "testing"

func TestCrackPassword(t *testing.T) {
	type args struct {
		archivePath string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
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
