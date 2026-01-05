package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	var downloadedExtractFiles []string
	for tag := range m {
		for _, asset := range m[tag] {
			assetName := asset.GetName()
			if c.assetTag {
				assetName = fileRename(assetName, tag)
			}
			if c.assetDate {
				updatedAt := asset.GetUpdatedAt()
				if !updatedAt.IsZero() {
					assetName = fileRename(assetName, updatedAt.Format("200601021504"))
				}
			}
			assetURL := asset.GetBrowserDownloadURL()
			assetPath := filepath.Join(c.path, assetName)

			if err := download(assetURL, assetPath); err != nil {
				return err
			}
			downloadedFiles = append(downloadedFiles, assetPath)

			if c.assetExtract {
				extractPath, err := extract(assetPath)
				if err != nil {
					return err
				}
				if extractPath != "" {
					downloadedExtractFiles = append(downloadedExtractFiles, extractPath)
				}
			}
		}
	}

	if c.autoclean {
		if err := cleanOldEntries(downloadedFiles, false); err != nil {
			log.Printf("failed to clean old release files: %v", err)
		}
		if err := cleanOldEntries(downloadedExtractFiles, true); err != nil {
			log.Printf("failed to clean old extracted directories: %v", err)
		}
	}

	return nil
}

func cleanOldEntries(newEntries []string, isDir bool) error {
	newEntriesSet := make(map[string]struct{}, len(newEntries))
	for _, f := range newEntries {
		newEntriesSet[f] = struct{}{}
	}

	entries, err := os.ReadDir(c.path)
	if err != nil {
		return err
	}

	entryType := "file"
	if isDir {
		entryType = "directory"
	}

	for _, entry := range entries {
		if entry.IsDir() != isDir {
			continue
		}

		entryPath := filepath.Join(c.path, entry.Name())

		if _, isNew := newEntriesSet[entryPath]; isNew {
			continue
		}

		log.Printf("removing old release %s: %s", entryType, entryPath)
		var removeErr error
		if isDir {
			removeErr = os.RemoveAll(entryPath)
		} else {
			removeErr = os.Remove(entryPath)
		}
		if removeErr != nil {
			log.Printf("failed to remove old release %s %s: %v", entryType, entryPath, removeErr)
		}
	}

	return nil
}
