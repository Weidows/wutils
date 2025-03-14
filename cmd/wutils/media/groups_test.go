package media

import "testing"

func TestClusterAndCopy(t *testing.T) {
	type args struct {
		inputDir string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test image",
			args: args{
				inputDir: "./test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ClusterAndCopy(tt.args.inputDir)
		})
	}
}
