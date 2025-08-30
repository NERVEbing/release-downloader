package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/google/go-github/v62/github"
)

const (
	DefaultTimeout  = time.Second * 30
	DefaultInterval = time.Hour
)

type config struct {
	repository   string
	proxy        string
	tag          string
	filename     string
	latest       bool
	prerelease   bool
	token        string
	path         string
	interval     time.Duration
	now          bool
	once         bool
	timeout      time.Duration
	assetTag     bool
	assetDate    bool
	assetExtract bool
	autoclean    bool

	httpClient   *http.Client
	githubClient *github.Client
}

var c *config

func init() {
	repositoryFlag := flag.String("repository", "", "GitHub repository in format `owner/repo`.")
	proxyFlag := flag.String("proxy", "", "Optional. HTTP/HTTPS proxy for downloads (e.g., `http://localhost:8080` or `socks5://127.0.0.1:1080`).")
	tagFlag := flag.String("tag", "", "Optional. Download assets matching a specific tag (supports regex, e.g., `.*.18.*`).")
	filenameFlag := flag.String("filename", "", "Optional. Filter assets by filename (regex supported, e.g., `.*linux-arm64.*.gz`).")
	latestFlag := flag.Bool("latest", false, "Optional. Fetch the latest release (ignores `tag` if set).")
	prereleaseFlag := flag.Bool("prerelease", false, "Optional. Include prerelease versions.")
	tokenFlag := flag.String("token", "", "Optional. GitHub personal access token (for private repos or higher rate limits).")
	pathFlag := flag.String("path", "./tmp", "Optional. Directory to save downloaded files (default: `./tmp`).")
	intervalFlag := flag.Duration("interval", DefaultInterval, "Optional. Delay between download tasks (e.g., `30s`, `5m`).")
	nowFlag := flag.Bool("now", false, "Optional. Run immediately without waiting for the first interval.")
	onceFlag := flag.Bool("once", false, "Optional. Run once and exit.")
	timeoutFlag := flag.Duration("timeout", DefaultTimeout, "Optional. HTTP client timeout (e.g., `30s`, `2m`).")
	assetTagFlag := flag.Bool("asset_tag", false, "Optional. Append release tag to filename (e.g., `file.zip` → `file-v1.0.0.zip`).")
	assetDateFlag := flag.Bool("asset_date", false, "Optional. Append download date to filename (e.g., `file.zip` → `file-20240502.zip`).")
	assetExtractFlag := flag.Bool("asset_extract", false, "Optional. Auto-extract downloaded files (supports `.zip`, `.gz`, `.tar.gz`).")
	autocleanFlag := flag.Bool("autoclean", false, "Optional. Remove old release files after downloading a new one.")

	flag.Parse()

	repository := envOrFlag("RD_REPOSITORY", *repositoryFlag)
	proxy := envOrFlag("RD_PROXY", *proxyFlag)
	tag := envOrFlag("RD_TAG", *tagFlag)
	filename := envOrFlag("RD_FILENAME", *filenameFlag)
	latest := envOrFlag("RD_LATEST", *latestFlag)
	prerelease := envOrFlag("RD_PRERELEASE", *prereleaseFlag)
	token := envOrFlag("RD_TOKEN", *tokenFlag)
	path := envOrFlag("RD_PATH", *pathFlag)
	interval := envOrFlag("RD_INTERVAL", *intervalFlag)
	now := envOrFlag("RD_NOW", *nowFlag)
	once := envOrFlag("RD_ONCE", *onceFlag)
	timeout := envOrFlag("RD_TIMEOUT", *timeoutFlag)
	assetTag := envOrFlag("RD_ASSET_TAG", *assetTagFlag)
	assetDate := envOrFlag("RD_ASSET_DATE", *assetDateFlag)
	assetExtract := envOrFlag("RD_ASSET_EXTRACT", *assetExtractFlag)
	autoclean := envOrFlag("RD_AUTOCLEAN", *autocleanFlag)

	httpClient := &http.Client{}
	if proxy != "" {
		proxy, err := url.Parse(proxy)
		if err != nil {
			log.Printf("invalid proxy URL: %s", proxy)
		} else {
			httpClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxy),
			}
		}
	}

	githubClient := github.NewClient(httpClient)

	if len(token) > 0 {
		githubClient = githubClient.WithAuthToken(token)
	}
	if timeout.Milliseconds() > 0 {
		httpClient.Timeout = timeout
	}

	c = &config{
		repository:   repository,
		proxy:        proxy,
		tag:          tag,
		filename:     filename,
		latest:       latest,
		prerelease:   prerelease,
		token:        token,
		path:         path,
		interval:     interval,
		now:          now,
		once:         once,
		timeout:      timeout,
		assetTag:     assetTag,
		assetDate:    assetDate,
		assetExtract: assetExtract,
		autoclean:    autoclean,

		httpClient:   httpClient,
		githubClient: githubClient,
	}
}

func main() {
	log.Printf("repository: %s", c.repository)
	log.Printf("proxy: %s", c.proxy)
	log.Printf("tag: %s", c.tag)
	log.Printf("filename: %s", c.filename)
	log.Printf("latest: %t", c.latest)
	log.Printf("prerelease: %t", c.prerelease)
	log.Printf("token: %s", c.token)
	log.Printf("path: %s", c.path)
	log.Printf("interval: %s", c.interval.String())
	log.Printf("now: %t", c.now)
	log.Printf("once: %t", c.once)
	log.Printf("timeout: %s", c.timeout.String())
	log.Printf("asset_tag: %t", c.assetTag)
	log.Printf("asset_date: %t", c.assetDate)
	log.Printf("asset_extract: %t", c.assetExtract)
	log.Printf("autoclean: %t", c.autoclean)

	ctx := context.Background()
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()

	if c.now {
		run(ctx)
		if c.once {
			return
		}
	}

	for range ticker.C {
		run(ctx)
		if c.once {
			return
		}
	}
}

func run(ctx context.Context) {
	t := time.Now()
	log.Printf("task started")

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

	log.Printf("task finished, duration: %s", time.Since(t))
}
