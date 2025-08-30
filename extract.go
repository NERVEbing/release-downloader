package main

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func extract(path string) (string, error) {
	name, ext := fileNameAndExt(path)
	target := name

	switch ext {
	case ".gz":
		target += "-" + "gz"
	case ".tar.gz":
		target += "-" + "tar-gz"
	case ".zip":
		target += "-" + "zip"
	default:
		log.Printf("unsupported file extension: %s, skip extract", ext)
		return "", nil
	}

	exist, err := isExistDir(target)
	if err != nil {
		return "", err
	}
	if exist {
		log.Printf("dir %s already exists, skip extract", target)
		return target, nil
	}

	log.Printf("extracting %s to %s", path, target)

	switch ext {
	case ".gz":
		err = extractGz(path, target)
	case ".tar.gz":
		err = extractTarGz(path, target)
	case ".zip":
		err = extractZip(path, target)
	}

	if err != nil {
		return "", err
	}

	log.Printf("extraction succeeded: %s", path)

	return target, nil
}

func extractGz(path string, target string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func() { _ = gzReader.Close() }()

	if err = os.MkdirAll(target, os.ModePerm); err != nil {
		return err
	}

	return fileWrite(filepath.Join(target, filepath.Base(target)), gzReader)
}

func extractTarGz(path string, target string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer func() { _ = gzReader.Close() }()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if strings.Contains(header.Name, "..") {
			log.Printf("skipping %s", header.Name)
			continue
		}

		outPath := filepath.Join(target, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if err = os.MkdirAll(outPath, os.ModePerm); err != nil {
				return err
			}
		case tar.TypeReg:
			if err = os.MkdirAll(filepath.Dir(outPath), os.ModePerm); err != nil {
				return err
			}
			if err = fileWrite(outPath, tarReader); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported file type %v", header.Typeflag)
		}
	}

	return nil
}

func extractZip(path string, target string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	rc, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()

	for _, f := range rc.File {
		if strings.Contains(f.Name, "..") {
			log.Printf("skipping %s", f.Name)
			continue
		}
		if err = extractZipFile(f, target); err != nil {
			return err
		}
	}

	return nil
}

func extractZipFile(f *zip.File, target string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer func() { _ = rc.Close() }()

	targetPath := filepath.Join(target, f.Name)
	if f.FileInfo().IsDir() {
		if err = os.MkdirAll(targetPath, os.ModePerm); err != nil {
			return err
		}
		return nil
	}
	if err = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
		return err
	}

	return fileWrite(targetPath, rc)
}
