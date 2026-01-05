## release-downloader

![Last Commit](https://custom-icon-badges.herokuapp.com/github/last-commit/NERVEbing/release-downloader?logo=history&logoColor=white)
![Build Docker Images](https://github.com/NERVEbing/release-downloader/actions/workflows/docker.yml/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/NERVEbing/release-downloader)
![License](https://custom-icon-badges.herokuapp.com/github/license/NERVEbing/release-downloader?logo=law&color=blue)

[English](README.md) | [简体中文](README.zh.md)

### Overview

release-downloader is an elegantly practical application designed to periodically monitor specified GitHub repositories
for releases and downloads corresponding files according to predefined criteria. It is primarily used to automate the
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

### English Version
| Command Line Flag | Environment Variable | Default Value     | Description                                                                                  |
| ----------------- | -------------------- | ----------------- | -------------------------------------------------------------------------------------------- |
| `-repository`     | `RD_REPOSITORY`      | `""`              | GitHub repository in format `owner/repo`.                                                    |
| `-proxy`          | `RD_PROXY`           | `""`              | HTTP/HTTPS proxy for downloads (e.g., `http://localhost:8080` or `socks5://127.0.0.1:1080`). |
| `-tag`            | `RD_TAG`             | `""`              | Download assets matching a specific tag (supports regex, e.g., `.*.18.*`).                   |
| `-filename`       | `RD_FILENAME`        | `""`              | Filter assets by filename (regex supported, e.g., `.*linux-arm64.*.gz`).                     |
| `-latest`         | `RD_LATEST`          | `false`           | Fetch the latest release (ignores `tag` if set).                                             |
| `-prerelease`     | `RD_PRERELEASE`      | `false`           | Include prerelease versions.                                                                 |
| `-token`          | `RD_TOKEN`           | `""`              | GitHub personal access token (for private repos or higher rate limits).                      |
| `-path`           | `RD_PATH`            | `"./tmp"`         | Directory to save downloaded files.                                                          |
| `-interval`       | `RD_INTERVAL`        | `DefaultInterval` | Delay between download tasks (e.g., `30s`, `5m`).                                            |
| `-now`            | `RD_NOW`             | `false`           | Run immediately without waiting for the first interval.                                      |
| `-once`           | `RD_ONCE`            | `false`           | Run once and exit.                                                                           |
| `-timeout`        | `RD_TIMEOUT`         | `DefaultTimeout`  | HTTP client timeout (e.g., `30s`, `2m`).                                                     |
| `-asset_tag`      | `RD_ASSET_TAG`       | `false`           | Append release tag to filename (e.g., `file.zip` → `file-v1.0.0.zip`).                       |
| `-asset_date`     | `RD_ASSET_DATE`      | `false`           | Append download date to filename (e.g., `file.zip` → `file-20240502.zip`).                   |
| `-asset_extract`  | `RD_ASSET_EXTRACT`   | `false`           | Auto-extract downloaded files (supports `.zip`, `.gz`, `.tar.gz`).                           |
| `-autoclean`      | `RD_AUTOCLEAN`       | `false`           | Remove old release files after downloading a new one.                                        |

### License

[Apache-2.0](LICENSE)
