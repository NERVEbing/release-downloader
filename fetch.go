package main

import (
	"context"
	"errors"
	"fmt"
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
			if c.assetTag {
				assetName = fileRename(assetName, tag)
			}
			if c.assetDate {
				assetName = fileRename(assetName, asset.UpdatedAt.Format("200601021504"))
			}
			assetURL := asset.GetBrowserDownloadURL()
			assetPath := filepath.Join(c.path, assetName)

			if err := download(assetURL, assetPath); err != nil {
				return err
			}

			if c.assetExtract {
				if err := extract(assetPath); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
