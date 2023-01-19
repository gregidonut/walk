package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		root string
		cfg  config
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{
			name: "NoFilter",
			args: args{
				root: "testdata",
				cfg: config{
					ext:  "",
					size: 0,
					list: true,
				},
			},
			wantOut: "testdata/dir.log\ntestdata/dir2/script.sh\n",
		},
		{
			name: "FilterExtensionMatch",
			args: args{
				root: "testdata",
				cfg: config{
					ext:  ".log",
					size: 0,
					list: true,
				},
			},
			wantOut: "testdata/dir.log\n",
		},
		{
			name: "FiterSizeMatch",
			args: args{
				root: "testdata",
				cfg: config{
					ext:  ".log",
					size: 10,
					list: true,
				},
			},
			wantOut: "testdata/dir.log\n",
		},
		{
			name: "FillerExtensionSizeNoMatch",
			args: args{
				root: "testdata",
				cfg: config{
					ext:  ".log",
					size: 20,
					list: true,
				},
			},
			wantOut: "",
		},
		{
			name: "FilterExtensionNoMatch",
			args: args{
				root: "testdata",
				cfg: config{
					ext:  ".gz",
					size: 0,
					list: true,
				},
			},
			wantOut: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := &bytes.Buffer{}
			err := run(tt.args.root, out, tt.args.cfg)
			if err != nil {
				t.Fatal(err)
			}

			if gotOut := out.String(); gotOut != tt.wantOut {
				t.Errorf("run() gotOut = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

// this helper will be used by a test function that will check if -del flag
// has deleted the files then delete the temp directory the files were
// created on
func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()

	tempDir, err := os.MkdirTemp("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fName := fmt.Sprintf("file%d%s", j, k)
			fPath := filepath.Join(tempDir, fName)
			if err := os.WriteFile(fPath, []byte("dummy"), 0643); err != nil {
				t.Fatal(err)
			}
		}
	}

	return tempDir, func() {
		os.RemoveAll(tempDir)
	}
}
