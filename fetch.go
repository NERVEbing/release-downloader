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

func fetchReleases(ctx context.Context) ([]*github.RepositoryRelease, error) {
	var r []*github.RepositoryRelease

	arr := strings.Split(c.repository, "/")
	if len(arr) != 2 {
		return nil, errors.New(fmt.Sprintf("invalid repository: %s", c.repository))
	}
	owner := arr[0]
	repo := arr[1]

	releases, _, err := c.githubClient.Repositories.ListReleases(ctx, owner, repo, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("list releases error: %s", err))
	}
	for _, release := range releases {
		if !c.prerelease && release.GetPrerelease() {
			continue
		}

		if c.latest {
			r = append(r, release)
			return r, nil
		}

		if len(c.tag) > 0 {
			match, err := matchPattern(release.GetTagName(), c.tag)
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

func fetchAssets(releases []*github.RepositoryRelease) (map[string][]*github.ReleaseAsset, error) {
	m := make(map[string][]*github.ReleaseAsset)

	if len(releases) == 0 {
		return m, nil
	}

	for _, release := range releases {
		for _, asset := range release.Assets {
			if len(c.filename) > 0 {
				match, err := matchPattern(asset.GetName(), c.filename)
				if err != nil {
					return nil, err
				}
				if !match {
					continue
				}
			}
			_, exist := m[release.GetTagName()]
			if exist {
				m[release.GetTagName()] = append(m[release.GetTagName()], asset)
			} else {
				m[release.GetTagName()] = []*github.ReleaseAsset{asset}
			}
		}
	}

	return m, nil
}

func fetchFiles(m map[string][]*github.ReleaseAsset) error {
	if len(m) == 0 {
		return nil
	}

	if _, err := os.Stat(c.path); err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err = os.MkdirAll(c.path, os.ModePerm); err != nil {
			return err
		}
	}

	for tag := range m {
		for _, asset := range m[tag] {
			assetName := asset.GetName()
			assetURL := asset.GetBrowserDownloadURL()
			assetTag := tag
			assetDate := asset.UpdatedAt.Format("200601021504")

			if c.assetTag {
				assetName = rename(assetName, assetTag)
			}
			if c.assetDate {
				assetName = rename(assetName, assetDate)
			}

			if err := download(assetName, assetURL, filepath.Join(c.path, assetName)); err != nil {
				if os.IsExist(err) {
					log.Printf("file %s already exists, skip", filepath.Join(c.path, assetName))
					continue
				}
				return err
			}
		}
	}

	return nil
}

func download(assetName string, assetURL string, assetPath string) error {
	exist, err := isExist(assetPath)
	if err != nil {
		return err
	}
	if exist {
		return os.ErrExist
	}

	resp, err := http.Get(assetURL)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("download status code %d", resp.StatusCode))
	}

	file, err := os.Create(assetPath)
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	if _, err = io.Copy(file, resp.Body); err != nil {
		return err
	}

	log.Printf("downloaded %s to %s success", assetName, assetPath)

	return nil
}
