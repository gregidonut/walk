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

func TestRunDelExtension(t *testing.T) {
	tests := []struct {
		name        string
		cfg         config
		extNoDelete string
		nDelete     int
		nNoDelete   int
		want        string
	}{
		{
			name: "DeleteExtensionNoMatch",
			cfg: config{
				ext: ".log",
				del: true,
			},
			extNoDelete: ".gz",
			nDelete:     0,
			nNoDelete:   10,
			want:        "",
		},
		{
			name: "DeleteExtensionMatch",
			cfg: config{
				ext: ".log",
				del: true,
			},
			extNoDelete: "",
			nDelete:     10,
			nNoDelete:   0,
			want:        "",
		},
		{
			name: "DeleteExtensionMixed",
			cfg: config{
				ext: ".log",
				del: true,
			},
			extNoDelete: ".gz",
			nDelete:     0,
			nNoDelete:   10,
			want:        "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}

			tempDir, cleanup := createTempDir(t, map[string]int{
				tt.cfg.ext:     tt.nDelete,
				tt.extNoDelete: tt.nNoDelete,
			})
			defer cleanup()

			if err := run(tempDir, w, tt.cfg); err != nil {
				t.Fatal(err)
			}

			if w.String() != tt.want {
				t.Errorf("want %q, got %q\n", tt.want, w.String())
			}

			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != tt.nNoDelete {
				t.Errorf("want %d files left, got %d instead\n", tt.nNoDelete, len(filesLeft))
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
