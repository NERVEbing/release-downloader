package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func fileWrite(path string, reader io.Reader) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()

	_, err = io.Copy(f, reader)
	return err
}

func fileNameAndExt(path string) (string, string) {
	ext := filepath.Ext(path)
	if strings.HasSuffix(strings.ToLower(path), ".tar.gz") {
		ext = ".tar.gz"
	}
	path = strings.TrimSuffix(path, ext)

	return path, ext
}

func fileRename(path string, opt ...string) string {
	name, ext := fileNameAndExt(path)
	for _, o := range opt {
		name += "-" + o
	}

	return name + ext
}

func isExistFile(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return !f.IsDir(), nil
}

func isExistDir(path string) (bool, error) {
	f, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return f.IsDir(), nil
}
