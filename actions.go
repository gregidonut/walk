package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func filterOut(path, ext string, minSize int64, info os.FileInfo) bool {
	// sad paths return true to simulate a continue statement in a forloop
	// that is checked by the calling function Run()
	// in other words if these conditions are met
	// i.e. path is a directory or size is less than minimum size specified
	// or extension specified does not match the file extension we will
	// *continue* otherwise we will return false and the calling function
	// in Run() will proceed with the next step which is to presumably list
	// the file if this function returns false
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	return false
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}
