## release-downloader

![Last Commit](https://custom-icon-badges.herokuapp.com/github/last-commit/NERVEbing/release-downloader?logo=history&logoColor=white)
![Build Docker Images](https://github.com/NERVEbing/release-downloader/actions/workflows/docker.yml/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/NERVEbing/release-downloader)
![License](https://custom-icon-badges.herokuapp.com/github/license/NERVEbing/release-downloader?logo=law&color=blue)

[English](README.md) | [简体中文](README.zh.md)

### Overview

release-downloader is an elegantly practical application designed to periodically monitor specified GitHub repositories
for releases and download corresponding files according to predefined criteria. It is primarily used to automate the
acquisition of the latest software release versions.

### Features

- Periodically monitors GitHub repositories for new releases.
- Supports matching based on release or tag names.
- Allows file name-based matching.
- Offers custom download paths.
- Automatically extracts files (supports zip, gz, and tar.gz).
- Developed in Go, ensuring simplicity and ease of use.

### Usage

#### Docker Compose (Recommended)

```shell
mkdir release-downloader && cd release-downloader
wget https://raw.githubusercontent.com/NERVEbing/release-downloader/master/docker-compose.yml
# Edit docker-compose.yml file
docker compose up -d
```

See [docker-compose.yml](docker-compose.yml)

#### Build from source

```shell
git clone https://github.com/NERVEbing/release-downloader.git && cd release-downloader
go build -o release-downloader .
./release-downloader \
    -repository "containers/podman" \
    -filename ".*linux_amd64.tar.gz$" \
    -latest \
    -prerelease \
    -path ./tmp \
    -interval 1h \
    -now
```

#### Environment & Command line arguments

| Command Line Flag | Environment Variable | Description                                                                                                        | Default Value |
| ----------------- | -------------------- | ------------------------------------------------------------------------------------------------------------------ | ------------- |
| `-repository`     | `RD_REPOSITORY`      | GitHub repository in {owner}/{repo} format.                                                                        | ""            |
| `-proxy`          | `RD_PROXY`           | Optional. HTTP/SOCKS proxy to use for downloads (e.g., 'http://localhost:8080' or 'socks5://127.0.0.1:1080').      | ""            |
| `-tag`            | `RD_TAG`             | Optional. Download assets from a specific tag (e.g., regexp '.*.18.*').                                            | ""            |
| `-filename`       | `RD_FILENAME`        | Optional. Download assets matching a specific filename, excluding tarball or zipball (e.g., '.*linux-arm64.*.gz'). | ""            |
| `-latest`         | `RD_LATEST`          | Optional. Download the latest release.                                                                             | `false`       |
| `-prerelease`     | `RD_PRERELEASE`      | Optional. Download prerelease versions.                                                                            | `false`       |
| `-token`          | `RD_TOKEN`           | Optional. GitHub personal access token.                                                                            | ""            |
| `-path`           | `RD_PATH`            | Optional. Download path.                                                                                           | "./tmp"       |
| `-interval`       | `RD_INTERVAL`        | Optional. Interval between download tasks.                                                                         | `1h`          |
| `-now`            | `RD_NOW`             | Optional. Run the task immediately.                                                                                | `false`       |
| `-timeout`        | `RD_TIMEOUT`         | Optional. HTTP client timeout duration.                                                                            | `30s`         |
| `-asset_tag`      | `RD_ASSET_TAG`       | Optional. Rename the file using the asset tag (e.g., xxx.zip -> xxx-v0.5.1.zip).                                   | `false`       |
| `-asset_date`     | `RD_ASSET_DATE`      | Optional. Rename the file using the asset date (e.g., xxx.zip -> xxx-202401231619.zip).                            | `false`       |
| `-asset_extract`  | `RD_ASSET_EXTRACT`   | Optional. Automatically extract files (supports zip, gz, and tar.gz).                                              | `false`       |

### License

[Apache-2.0](LICENSE)
