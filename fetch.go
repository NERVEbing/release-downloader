package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-github/v62/github"
)

func fetchReleases(ctx context.Context, client *github.Client, repository string, tag string, latest bool, prerelease bool) ([]*github.RepositoryRelease, error) {
	var r []*github.RepositoryRelease

	arr := strings.Split(repository, "/")
	if len(arr) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid repository: %s", repository))
	}
	owner := arr[0]
	repo := arr[1]

	releases, _, err := client.Repositories.ListReleases(ctx, owner, repo, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("list releases error: %s", err))
	}
	for _, release := range releases {
		if !prerelease && release.GetPrerelease() {
			continue
		}

		if latest {
			r = append(r, release)
			return r, nil
		}

		if len(tag) > 0 {
			match, err := matchPattern(release.GetTagName(), tag)
			if err != nil {
				return nil, err
			}
			if !match {
				continue
			}
		}

		r = append(r, release)
	}

	return r, nil
}

func fetchAssets(releases []*github.RepositoryRelease, filename string) ([]*github.ReleaseAsset, error) {
	var a []*github.ReleaseAsset

	if len(releases) == 0 {
		return a, nil
	}

	for _, release := range releases {
		for _, asset := range release.Assets {
			if len(filename) > 0 {
				match, err := matchPattern(asset.GetName(), filename)
				if err != nil {
					return nil, err
				}
				if !match {
					continue
				}
			}
			a = append(a, asset)
		}
	}

	return a, nil
}

func fetchFiles(m map[string]string, path string) error {
	if len(m) == 0 {
		return nil
	}

	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err = os.MkdirAll(path, os.ModePerm); err != nil {
			return err
		}
	}

	for n, u := range m {
		if err := download(n, u, path); err != nil {
			if os.IsExist(err) {
				log.Printf("file %s already exists, skip", filepath.Join(path, n))
				continue
			}
			return err
		}
	}

	return nil

}

func download(n string, u string, path string) error {
	p := filepath.Join(path, n)
	f, err := os.Stat(p)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		return os.ErrExist
	}
	if f != nil && f.Size() < 1024 {
		if err = os.Remove(p); err != nil {
			return err
		}
	}

	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("download status code %d", resp.StatusCode))
	}

	file, err := os.Create(p)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}

	log.Printf("downloaded %s to %s success", n, p)

	return nil
}
