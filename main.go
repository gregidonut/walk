package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	// extension to filter out
	ext []string
	// min file size
	size int64
	// length of file name
	nameLength int
	// list files
	list bool
	// delete files
	del bool
	// log destination writer
	wLog io.Writer
	// archive directory
	archive string
}

func main() {
	// Parsing command line flags
	root := flag.String("root", ".", "Root directory to start")
	logFile := flag.String("log", "", "Log deletes to this file")
	// Action options
	list := flag.Bool("list", false, "List files only")
	del := flag.Bool("del", false, "Delete files")
	archive := flag.String("archive", "", "Archive directory")
	// Filter options
	ext := flag.Bool("ext", false, "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	nameLength := flag.Int("name-length", 0, "minimum characters in file name")
	flag.Parse()

	var (
		f   = os.Stdout
		err error
	)
	if *logFile != "" {
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
	}

	extensions := make([]string, 0)
	if *ext {
		extensions = append(extensions, flag.Args()...)
	}

	c := config{
		ext:        extensions,
		nameLength: *nameLength,
		size:       *size,
		list:       *list,
		del:        *del,
		wLog:       f,
		archive:    *archive,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	delLogger := log.New(cfg.wLog, "DELETED FILE: ", log.LstdFlags)
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filterOut(path, cfg.ext, cfg.size, cfg.nameLength, info) {
			return nil
		}

		// If list was explicitly set, don't do anything else
		if cfg.list {
			return listFile(path, out)
		}

		// Archive files and continue if successful
		if cfg.archive != "" {
			if err := archiveFile(cfg.archive, root, path); err != nil {
				return err
			}
		}

		if cfg.del {
			return delFile(path, delLogger)
		}

		// List is the default option if nothing else was set
		return listFile(path, out)
	})
}
