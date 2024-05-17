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
```

#### Environment & Command line arguments

| Command Line Flag | Environment Variable | Description                                                                                                        | Default Value |
|-------------------|----------------------|--------------------------------------------------------------------------------------------------------------------|---------------|
| `-repository`     | `RD_REPOSITORY`      | GitHub repository in {owner}/{repo} format.                                                                        | ""            |
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

### License

[Apache-2.0](LICENSE)
