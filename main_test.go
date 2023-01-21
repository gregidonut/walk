package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
					ext:  make([]string, 0),
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
					ext:  []string{".log"},
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
					ext:  []string{".log"},
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
					ext:  []string{".log"},
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
					ext:  []string{".gz"},
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
				ext: []string{".log"},
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
				ext: []string{".log"},
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
				ext: []string{".log"},
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
			buffer := &bytes.Buffer{}
			logBuffer := &bytes.Buffer{}
			tt.cfg.wLog = logBuffer

			filesMap := make(map[string]int, 0)
			for _, e := range tt.cfg.ext {
				filesMap[e] = tt.nDelete
			}
			filesMap[tt.extNoDelete] = tt.nNoDelete

			tempDir, cleanup := createTempDir(t, filesMap)
			defer cleanup()

			if err := run(tempDir, buffer, tt.cfg); err != nil {
				t.Fatal(err)
			}

			if buffer.String() != tt.want {
				t.Errorf("want %q, got %q\n", tt.want, buffer.String())
			}

			filesLeft, err := os.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			if len(filesLeft) != tt.nNoDelete {
				t.Errorf("want %d files left, got %d instead\n", tt.nNoDelete, len(filesLeft))
			}

			want := tt.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != want {
				t.Errorf("want %d log lines, got %d \n", want, len(lines))
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

func TestRunArchive(t *testing.T) {
	tests := []struct {
		name         string
		cfg          config
		extNoArchive string
		nArchive     int
		nNoArchive   int
	}{
		{
			name: "ArchiveExtensionNoMatch",
			cfg: config{
				ext: []string{".log"},
			},
			extNoArchive: ".gz",
			nArchive:     0,
			nNoArchive:   10,
		},
		{
			name: "ArchiveExtensionMatch",
			cfg: config{
				ext: []string{".log"},
			},
			extNoArchive: "",
			nArchive:     10,
			nNoArchive:   0,
		},
		{
			name: "ArchiveExtensionMixed",
			cfg: config{
				ext: []string{".log"},
			},
			extNoArchive: ".gz",
			nArchive:     5,
			nNoArchive:   5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			filesMap := make(map[string]int, 0)
			for _, e := range tt.cfg.ext {
				filesMap[e] = tt.nArchive
			}
			filesMap[tt.extNoArchive] = tt.nNoArchive

			tempDir, cleanup := createTempDir(t, filesMap)
			defer cleanup()

			archiveDir, cleanupArchive := createTempDir(t, nil)
			defer cleanupArchive()

			tt.cfg.archive = archiveDir
			if err := run(tempDir, buffer, tt.cfg); err != nil {
				t.Fatal(err)
			}

			pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", tt.cfg.ext))
			expFiles, err := filepath.Glob(pattern)
			if err != nil {
				t.Fatal(err)
			}

			want := strings.Join(expFiles, "\n")
			got := strings.TrimSpace(buffer.String())
			if got != want {
				t.Errorf("want %q, got %q\n", want, got)
			}

			filesArchived, err := os.ReadDir(archiveDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(filesArchived) != tt.nArchive {
				t.Errorf("want %d files archived, got %d", tt.nArchive, len(filesArchived))
			}
		})
	}
}
