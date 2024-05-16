## release-downloader

### Overview

release-downloader is an application designed to periodically monitor specified GitHub repositories for releases and
download corresponding files based on defined criteria. It is primarily used for automating the process of fetching the
latest software release versions.

### Features

- Periodically monitors GitHub repositories for releases.
- Supports matching based on release or tag names.
- Allows matching based on file names.
- Provides custom download paths.
- Developed in Go, offering simplicity and ease of use.

### Usage

#### Docker Compose(Recommended)

```shell
mkdir release-downloader && cd release-downloader
wget https://raw.githubusercontent.com/NERVEbing/release-downloader/master/docker-compose.yml
# Edit docker-compose.yml file
docker compose up -d
```

##### Environment

See [docker-compose.yml](docker-compose.yml)

#### Build from source

```shell
go build

./release-downloader \
    -repository "fastfetch-cli/fastfetch" \
    -filename ".*linux.*amd64.*.zip" \
    -latest \
    -prerelease \
    -path ./tmp \
    -interval 1h \
    -now

2024/05/16 15:26:27 repository: fastfetch-cli/fastfetch
2024/05/16 15:26:27 tag: 
2024/05/16 15:26:27 filename: .*linux.*amd64.*.zip
2024/05/16 15:26:27 latest: true
2024/05/16 15:26:27 prerelease: true
2024/05/16 15:26:27 token:
2024/05/16 15:26:27 path: ./tmp
2024/05/16 15:26:27 interval: 1h0m0s
2024/05/16 15:26:27 now: true
2024/05/16 15:26:27 task commencing
2024/05/16 15:26:28 downloaded fastfetch-linux-amd64.zip to tmp/fastfetch-linux-amd64.zip success
2024/05/16 15:26:28 task completed, duration: 3.270727875s
```

##### Command line arguments

```shell
./release-downloader -h
Usage of ./release-downloader:
  -filename string
        optional, default: '', download assets from a specific filename, excluding tarball or zipball (eg: regexp '.*linux-arm64.*.gz')
  -interval duration
        optional, default: 1h, download task interval (default 1h0m0s)
  -latest
        optional, default: false, download target as latest release
  -now
        optional, default: false, immediately run
  -path string
        optional, default: './tmp', download path (default "./tmp")
  -prerelease
        optional, default: false, download target as prerelease
  -repository string
        github repository, {owner}/{repo}
  -tag string
        optional, default: '', download assets from a specific tag (eg: regexp '.*.18.*')
  -token string
        optional, default: '', github personal access token

```

### License

[Apache-2.0](LICENSE)
