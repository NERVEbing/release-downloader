package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/google/go-github/v62/github"
)

func fetchReleases(ctx context.Context) ([]*github.RepositoryRelease, error) {
	var r []*github.RepositoryRelease

	arr := strings.Split(c.repository, "/")
	if len(arr) != 2 {
		return nil, fmt.Errorf("invalid repository: %s", c.repository)
	}
	owner := arr[0]
	repo := arr[1]

	releases, _, err := c.githubClient.Repositories.ListReleases(ctx, owner, repo, nil)
	if err != nil {
		return nil, fmt.Errorf("list releases error: %s", err)
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

	var downloadedFiles []string
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
			downloadedFiles = append(downloadedFiles, assetPath)
		}
	}

	if c.autoclean {
		if err := cleanOldFiles(c.path, downloadedFiles, c.filename); err != nil {
			log.Printf("failed to clean old releases: %v", err)
		}
	}

	return nil
}

func cleanOldFiles(downloadPath string, newFiles []string, filenameRegex string) error {
	if filenameRegex == "" {
		return nil
	}
	re, err := regexp.Compile(filenameRegex)
	if err != nil {
		return fmt.Errorf("invalid filename regex for cleaning: %w", err)
	}

	newFilesSet := make(map[string]struct{}, len(newFiles))
	for _, f := range newFiles {
		newFilesSet[f] = struct{}{}
	}

	files, err := os.ReadDir(downloadPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(downloadPath, file.Name())

		if _, isNew := newFilesSet[filePath]; isNew {
			continue
		}

		if re.MatchString(file.Name()) {
			log.Printf("removing old release file: %s", filePath)
			if err := os.Remove(filePath); err != nil {
				log.Printf("failed to remove old release file %s: %v", filePath, err)
			}
		}
	}

	return nil
}
