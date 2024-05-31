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
	DefaultTimeout  = time.Second * 30
	DefaultInterval = time.Hour
)

type config struct {
	repository   string
	tag          string
	filename     string
	latest       bool
	prerelease   bool
	token        string
	path         string
	interval     time.Duration
	now          bool
	timeout      time.Duration
	assetTag     bool
	assetDate    bool
	assetExtract bool

	httpClient   *http.Client
	githubClient *github.Client
}

var c *config

func init() {
	repositoryFlag := flag.String("repository", "", "GitHub repository in {owner}/{repo} format.")
	tagFlag := flag.String("tag", "", "Optional. Download assets from a specific tag (e.g., regexp '.*.18.*').")
	filenameFlag := flag.String("filename", "", "Optional. Download assets matching a specific filename, excluding tarball or zipball (e.g., '.*linux-arm64.*.gz').")
	latestFlag := flag.Bool("latest", false, "Optional. Download the latest release.")
	prereleaseFlag := flag.Bool("prerelease", false, "Optional. Download prerelease versions.")
	tokenFlag := flag.String("token", "", "Optional. GitHub personal access token.")
	pathFlag := flag.String("path", "./tmp", "Optional. Download path.")
	intervalFlag := flag.Duration("interval", DefaultInterval, "Optional. Interval between download tasks.")
	nowFlag := flag.Bool("now", false, "Optional. Run the task immediately.")
	timeoutFlag := flag.Duration("timeout", DefaultTimeout, "Optional. HTTP client timeout duration.")
	assetTagFlag := flag.Bool("asset_tag", false, "Optional. Rename the file using the asset tag (e.g., xxx.zip -> xxx-v0.5.1.zip).")
	assetDateFlag := flag.Bool("asset_date", false, "Optional. Rename the file using the asset date (e.g., xxx.zip -> xxx-202401231619.zip).")
	assetExtractFlag := flag.Bool("asset_extract", false, "Optional. Automatically extract files (supports zip, gz, and tar.gz).")
	flag.Parse()

	repository := envOrFlag("RD_REPOSITORY", *repositoryFlag)
	tag := envOrFlag("RD_TAG", *tagFlag)
	filename := envOrFlag("RD_FILENAME", *filenameFlag)
	latest := envOrFlag("RD_LATEST", *latestFlag)
	prerelease := envOrFlag("RD_PRERELEASE", *prereleaseFlag)
	token := envOrFlag("RD_TOKEN", *tokenFlag)
	path := envOrFlag("RD_PATH", *pathFlag)
	interval := envOrFlag("RD_INTERVAL", *intervalFlag)
	now := envOrFlag("RD_NOW", *nowFlag)
	timeout := envOrFlag("RD_TIMEOUT", *timeoutFlag)
	assetTag := envOrFlag("RD_ASSET_TAG", *assetTagFlag)
	assetDate := envOrFlag("RD_ASSET_DATE", *assetDateFlag)
	assetExtract := envOrFlag("RD_ASSET_EXTRACT", *assetExtractFlag)

	httpClient := &http.Client{}
	if timeout.Milliseconds() > 0 {
		httpClient.Timeout = timeout
	}
	githubClient := github.NewClient(httpClient)
	if len(token) > 0 {
		githubClient = githubClient.WithAuthToken(token)
	}

	c = &config{
		repository:   repository,
		tag:          tag,
		filename:     filename,
		latest:       latest,
		prerelease:   prerelease,
		token:        token,
		path:         path,
		interval:     interval,
		now:          now,
		timeout:      timeout,
		assetTag:     assetTag,
		assetDate:    assetDate,
		assetExtract: assetExtract,

		httpClient:   httpClient,
		githubClient: githubClient,
	}
}

func main() {
	log.Printf("repository: %s", c.repository)
	log.Printf("tag: %s", c.tag)
	log.Printf("filename: %s", c.filename)
	log.Printf("latest: %t", c.latest)
	log.Printf("prerelease: %t", c.prerelease)
	log.Printf("token: %s", c.token)
	log.Printf("path: %s", c.path)
	log.Printf("interval: %s", c.interval.String())
	log.Printf("now: %t", c.now)
	log.Printf("timeout: %s", c.timeout.String())
	log.Printf("asset_tag: %t", c.assetTag)
	log.Printf("asset_date: %t", c.assetDate)
	log.Printf("asset_extract: %t", c.assetExtract)

	ctx := context.Background()
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	if c.now {
		run(ctx)
	}

	for {
		select {
		case <-ticker.C:
			run(ctx)
		}
	}
}

func run(ctx context.Context) {
	t := time.Now()
	log.Printf("task commencing")

	releases, err := fetchReleases(ctx)
	if err != nil {
		log.Printf("failed to fetch releases: %v", err)
		return
	}

	m, err := fetchAssets(releases)
	if err != nil {
		log.Printf("failed to fetch assets: %v", err)
		return
	}

	if err = fetchFiles(m); err != nil {
		log.Printf("failed to fetch files: %v", err)
		return
	}

	log.Printf("task completed, duration: %s", time.Since(t))
}
