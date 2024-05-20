package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func download(link string, path string) error {
	exist, err := isExistFile(path)
	if err != nil {
		return err
	}
	if exist {
		log.Printf("file %s already exists, skip download", path)
		return nil
	}

	log.Printf("downloading %s", link)

	resp, err := http.Get(link)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("download status code %d", resp.StatusCode))
	}

	if err = fileWrite(path, resp.Body); err != nil {
		return err
	}

	log.Printf("download succeeded: %s", path)

	return nil
}
