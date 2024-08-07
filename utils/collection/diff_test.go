package collection

import (
	"reflect"
	"testing"
)

func TestSliceDiff(t *testing.T) {
	type args struct {
		a []string
		b []string
	}
	tests := []struct {
		name        string
		args        args
		wantMissInA []string
		wantMissInB []string
	}{
		{
			name: "No differences",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"a", "b", "c"},
			},
			wantMissInA: []string{},
			wantMissInB: []string{},
		},
		{
			name: "All elements different",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"d", "e", "f"},
			},
			wantMissInA: []string{"d", "e", "f"},
			wantMissInB: []string{"a", "b", "c"},
		},
		{
			name: "Some elements different",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{"b", "c", "d"},
			},
			wantMissInA: []string{"d"},
			wantMissInB: []string{"a"},
		},
		{
			name: "Empty slice a",
			args: args{
				a: []string{},
				b: []string{"a", "b", "c"},
			},
			wantMissInA: []string{"a", "b", "c"},
			wantMissInB: []string{},
		},
		{
			name: "Empty slice b",
			args: args{
				a: []string{"a", "b", "c"},
				b: []string{},
			},
			wantMissInA: []string{},
			wantMissInB: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMissInA, gotMissInB := SliceDiff(tt.args.a, tt.args.b)
			if !reflect.DeepEqual(gotMissInA, tt.wantMissInA) {
				t.Errorf("SliceDiff() gotMissInA = %v, want %v", gotMissInA, tt.wantMissInA)
			}
			if !reflect.DeepEqual(gotMissInB, tt.wantMissInB) {
				t.Errorf("SliceDiff() gotMissInB = %v, want %v", gotMissInB, tt.wantMissInB)
			}
		})
	}
}
