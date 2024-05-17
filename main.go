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

var (
	repositoryFlag = flag.String("repository", "", "github repository, {owner}/{repo}")
	tagFlag        = flag.String("tag", "", "optional, default: '', download assets from a specific tag (eg: regexp '.*.18.*')")
	filenameFlag   = flag.String("filename", "", "optional, default: '', download assets from a specific filename, excluding tarball or zipball (eg: regexp '.*linux-arm64.*.gz')")
	latestFlag     = flag.Bool("latest", false, "optional, default: false, download target as latest release")
	prereleaseFlag = flag.Bool("prerelease", false, "optional, default: false, download target as prerelease")
	tokenFlag      = flag.String("token", "", "optional, default: '', github personal access token")
	pathFlag       = flag.String("path", "./tmp", "optional, default: './tmp', download path")
	intervalFlag   = flag.Duration("interval", DefaultInterval, "optional, default: 1h, download task interval")
	nowFlag        = flag.Bool("now", false, "optional, default: false, immediately run")
	timeoutFlag    = flag.Duration("timeout", DefaultTimeout, "optional, default: 30s, http client timeout")
	assetTagFlag   = flag.Bool("asset_tag", false, "optional, default: false, rename the file using the asses tag (eg: xxx.zip -> xxx-v0.5.1.zip)")
	assetDateFlag  = flag.Bool("asset_date", false, "optional, default: false, rename the file using the asses date (eg: xxx.zip -> xxx-202401231619.zip)")
)

type config struct {
	repository string
	tag        string
	filename   string
	latest     bool
	prerelease bool
	token      string
	path       string
	interval   time.Duration
	now        bool
	timeout    time.Duration
	assetTag   bool
	assetDate  bool

	httpClient   *http.Client
	githubClient *github.Client
}

var c *config

func init() {
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

	httpClient := &http.Client{}
	if timeout.Milliseconds() > 0 {
		httpClient.Timeout = timeout
	}
	githubClient := github.NewClient(httpClient)
	if len(token) > 0 {
		githubClient = githubClient.WithAuthToken(token)
	}

	c = &config{
		repository: repository,
		tag:        tag,
		filename:   filename,
		latest:     latest,
		prerelease: prerelease,
		token:      token,
		path:       path,
		interval:   interval,
		now:        now,
		timeout:    timeout,
		assetTag:   assetTag,
		assetDate:  assetDate,

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
		log.Fatalf("failed to fetch releases: %v", err)
	}

	m, err := fetchAssets(releases)
	if err != nil {
		log.Fatalf("failed to fetch assets: %v", err)
	}

	if err = fetchFiles(m); err != nil {
		log.Fatalf("failed to fetch files: %v", err)
	}

	log.Printf("task completed, duration: %s", time.Since(t))
}
