package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/google/go-github/v62/github"
)

const (
	DefaultHttpTimeout = time.Second * 30
	DefaultInterval    = time.Hour
)

func main() {
	var (
		repository string
		tag        string
		filename   string
		latest     bool
		prerelease bool
		token      string
		path       string
		interval   time.Duration
		now        bool
	)

	flagRepository := flag.String("repository", "", "repository {owner}/{repo} (eg: fastfetch-cli/fastfetch)")
	flagTag := flag.String("tag", "", "download assets from a specific tag (eg: regexp '.*.18.*')")
	flagFilename := flag.String("filename", "", "download assets from a specific filename, excluding tarball or zipball (eg: regexp '.*linux-arm64.*.gz')")
	flagLatest := flag.Bool("latest", false, "download target as latest release")
	flagPrerelease := flag.Bool("prerelease", false, "download target as prerelease release")
	flagToken := flag.String("token", "", "github personal access token")
	flagPath := flag.String("path", "./tmp", "download path")
	flagInterval := flag.Duration("interval", DefaultInterval, "download interval")
	flogNow := flag.Bool("now", false, "immediately run")
	flag.Parse()

	repository = getEnvOrDefault("ENV_REPOSITORY", *flagRepository)
	tag = getEnvOrDefault("ENV_TAG", *flagTag)
	filename = getEnvOrDefault("ENV_FILE", *flagFilename)
	latest = getEnvOrDefaultBool("ENV_LATEST", *flagLatest)
	prerelease = getEnvOrDefaultBool("ENV_PRERELEASE", *flagPrerelease)
	token = getEnvOrDefault("ENV_TOKEN", *flagToken)
	path = getEnvOrDefault("ENV_PATH", *flagPath)
	interval = getEnvOrDefaultDuration("ENV_INTERVAL", *flagInterval)
	now = getEnvOrDefaultBool("ENV_NOW", *flogNow)

	log.Printf("repository: %s", repository)
	log.Printf("tag: %s", tag)
	log.Printf("filename: %s", filename)
	log.Printf("latest: %t", latest)
	log.Printf("prerelease: %t", prerelease)
	log.Printf("token: %s", token)
	log.Printf("path: %s", path)
	log.Printf("interval: %s", interval.String())
	log.Printf("now: %t", now)

	ctx := context.Background()
	client := github.NewClient(&http.Client{Timeout: DefaultHttpTimeout}).WithAuthToken(token)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	if now {
		run(ctx, client, repository, tag, filename, latest, prerelease, path)
	}

	for {
		select {
		case <-ticker.C:
			run(ctx, client, repository, tag, filename, latest, prerelease, path)
		}
	}
}

func run(ctx context.Context, client *github.Client, repository string, tag string, filename string, latest bool, prerelease bool, path string) {
	t := time.Now()
	log.Printf("task commencing")

	releases, err := fetchReleases(ctx, client, repository, tag, latest, prerelease)
	if err != nil {
		log.Fatalf("failed to fetch releases: %v", err)
	}

	assets, err := fetchAssets(releases, filename)
	if err != nil {
		log.Fatalf("failed to fetch assets: %v", err)
	}

	m := make(map[string]string)
	for _, asset := range assets {
		if _, exist := m[asset.GetName()]; !exist {
			m[asset.GetName()] = asset.GetBrowserDownloadURL()
		}
	}

	if err = fetchFiles(m, path); err != nil {
		log.Fatalf("failed to fetch files: %v", err)
	}

	log.Printf("task completed, duration: %s", time.Since(t))
}
