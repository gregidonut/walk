package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
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

func delFile(path string, delLogger *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	delLogger.Println(path)
	return nil
}

func archiveFile(destDir, root, path string) error {
	targetPath, err := relDirPath(destDir, root, path)
	if err != nil {
		return err
	}

	// making the directory structure in the specified destination
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	return writeToCompressedFile(targetPath, path)
}

func relDirPath(destDir, root, path string) (string, error) {
	// making sure no errors from checking specified destination dir info
	info, err := os.Stat(destDir)
	if err != nil {
		return "", err
	}

	// if destination specified is not a directory return error
	if !info.IsDir() {
		return "", fmt.Errorf("%s is not a directory", destDir)
	}

	// get the relative path from the root specified
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return "", err
	}

	// string representing the filename with added '.gz' at the end
	dest := fmt.Sprintf("%s.gz", filepath.Base(path))

	// a full path, with destination dir then relative path (to
	// preserve dir structure when it wasn't compressed) then the
	// filename + .gz
	return filepath.Join(destDir, relDir, dest), nil
}

func writeToCompressedFile(targetPath, inputPath string) error {
	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer out.Close()

	in, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer in.Close()

	zw := gzip.NewWriter(out)
	zw.Name = filepath.Base(inputPath)

	if _, err = io.Copy(zw, in); err != nil {
		return err
	}

	return out.Close()
}
