## release-downloader

![Last Commit](https://custom-icon-badges.herokuapp.com/github/last-commit/NERVEbing/release-downloader?logo=history&logoColor=white)
![Build Docker Images](https://github.com/NERVEbing/release-downloader/actions/workflows/docker.yml/badge.svg)
![Go Report Card](https://goreportcard.com/badge/github.com/NERVEbing/release-downloader)
![License](https://custom-icon-badges.herokuapp.com/github/license/NERVEbing/release-downloader?logo=law&color=blue)

[English](README.md) | [简体中文](README.zh.md)

### 概览

release-downloader 是一个简单实用的应用程序，旨在定期监视指定的 GitHub 仓库中的
releases，并根据预定标准下载相应的文件。它主要用于自动获取最新的软件发行版本。

### 主要功能

- 定期监控 GitHub 仓库中的新版本。
- 支持基于发布或标签名称的匹配。
- 允许基于文件名进行匹配。
- 提供自定义下载路径。
- 自动解压文件（支持 zip、gz 和 tar.gz）。
- 使用 Go 开发，确保简便易用。

### 使用方法

#### 使用 Docker Compose (推荐)

```shell
mkdir release-downloader && cd release-downloader
wget https://raw.githubusercontent.com/NERVEbing/release-downloader/master/docker-compose.yml
# 编辑 docker-compose.yml 文件
docker compose up -d
```

查看 [docker-compose.yml](docker-compose.yml)

#### 源码编译

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

#### 环境变量 & 命令行标志

| 命令行标志            | 环境变量               | 描述                                                                | 默认值     |
|------------------|--------------------|-------------------------------------------------------------------|---------|
| `-repository`    | `RD_REPOSITORY`    | 以 {owner}/{repo} 格式指定 GitHub 仓库。                                  | ""      |
| `-tag`           | `RD_TAG`           | 可选。从特定标签下载资产 (例如, regexp '.*.18.*')。                              | ""      |
| `-filename`      | `RD_FILENAME`      | 可选。下载与特定文件名匹配的资产，排除 tarball 或 zipball (例如, '.*linux-arm64.*.gz')。 | ""      |
| `-latest`        | `RD_LATEST`        | 可选。下载最新版本。                                                        | `false` |
| `-prerelease`    | `RD_PRERELEASE`    | 可选。下载预发行版本。                                                       | `false` |
| `-token`         | `RD_TOKEN`         | 可选。GitHub 个人访问令牌。                                                 | ""      |
| `-path`          | `RD_PATH`          | 可选。下载路径。                                                          | "./tmp" |
| `-interval`      | `RD_INTERVAL`      | 可选。下载任务之间的间隔时间。                                                   | `1h`    |
| `-now`           | `RD_NOW`           | 可选。立即运行任务。                                                        | `false` |
| `-timeout`       | `RD_TIMEOUT`       | 可选。HTTP 客户端超时时间。                                                  | `30s`   |
| `-asset_tag`     | `RD_ASSET_TAG`     | 可选。使用资产标签重命名文件 (例如, xxx.zip -> xxx-v0.5.1.zip)。                   | `false` |
| `-asset_date`    | `RD_ASSET_DATE`    | 可选。使用资产日期重命名文件 (例如, xxx.zip -> xxx-202401231619.zip)。             | `false` |
| `-asset_extract` | `RD_ASSET_EXTRACT` | 可选。自动提取文件 (支持 zip, gz, and tar.gz)。                               | `false` |

### 许可证

[Apache-2.0](LICENSE)
