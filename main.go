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

	httpClient   *http.Client
	githubClient *github.Client
}

var c *config

func init() {
	flag.Parse()

	repository := envOrFlag("REPOSITORY", *repositoryFlag)
	tag := envOrFlag("TAG", *tagFlag)
	filename := envOrFlag("FILENAME", *filenameFlag)
	latest := envOrFlag("LATEST", *latestFlag)
	prerelease := envOrFlag("PRERELEASE", *prereleaseFlag)
	token := envOrFlag("TOKEN", *tokenFlag)
	path := envOrFlag("PATH", *pathFlag)
	interval := envOrFlag("INTERVAL", *intervalFlag)
	now := envOrFlag("NOW", *nowFlag)

	httpClient := &http.Client{Timeout: DefaultHttpTimeout}
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

	assets, err := fetchAssets(releases)
	if err != nil {
		log.Fatalf("failed to fetch assets: %v", err)
	}

	if err = fetchFiles(assets); err != nil {
		log.Fatalf("failed to fetch files: %v", err)
	}

	log.Printf("task completed, duration: %s", time.Since(t))
}
