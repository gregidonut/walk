package main

import (
	"os"
	"testing"
)

func Test_filterOut(t *testing.T) {
	type args struct {
		path    string
		ext     string
		minSize int64
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "FilterNoExtension",
			args: args{
				path:    "testdata/dir.log",
				ext:     "",
				minSize: 0,
			},
			want: false,
		},
		{
			name: "FilterExtensionMatch",
			args: args{
				path:    "testdata/dir.log",
				ext:     ".log",
				minSize: 0,
			},
			want: false,
		},
		{
			name: "FilterExtensionNoMatch",
			args: args{
				path:    "testdata/dir.log",
				ext:     ".sh",
				minSize: 0,
			},
			want: true,
		},
		{
			name: "FilterExtensionSizeMatch",
			args: args{
				path:    "testdata/dir.log",
				ext:     ".log",
				minSize: 10,
			},
			want: false,
		},
		{
			name: "FilterExtensionSizeNoMatch",
			args: args{
				path:    "testdata/dir.log",
				ext:     ".log",
				minSize: 20,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := os.Stat(tt.args.path)
			if err != nil {
				t.Fatal(err)
			}

			got := filterOut(tt.args.path, tt.args.ext, tt.args.minSize, info)

			if got != tt.want {
				t.Errorf("filterOut() = %v, want %v", got, tt.want)
			}
		})
	}
}
