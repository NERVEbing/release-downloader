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

| 命令行参数       | 环境变量           | 默认值            | 说明                                                                                      |
| ---------------- | ------------------ | ----------------- | ----------------------------------------------------------------------------------------- |
| `-repository`    | `RD_REPOSITORY`    | `""`              | GitHub 仓库，格式为 `owner/repo`。                                                        |
| `-proxy`         | `RD_PROXY`         | `""`              | 下载时使用的 HTTP/HTTPS 代理（如 `http://localhost:8080` 或 `socks5://127.0.0.1:1080`）。 |
| `-tag`           | `RD_TAG`           | `""`              | 下载符合指定标签的资源（支持正则表达式，如 `.*.18.*`）。                                  |
| `-filename`      | `RD_FILENAME`      | `""`              | 按文件名筛选资源（支持正则，如 `.*linux-arm64.*.gz`）。                                   |
| `-latest`        | `RD_LATEST`        | `false`           | 下载最新版本（若设置，则忽略 `tag`）。                                                    |
| `-prerelease`    | `RD_PRERELEASE`    | `false`           | 包含预发布版本。                                                                          |
| `-token`         | `RD_TOKEN`         | `""`              | GitHub 个人访问令牌（用于私有仓库或提高速率限制）。                                       |
| `-path`          | `RD_PATH`          | `"./tmp"`         | 文件下载保存目录。                                                                        |
| `-interval`      | `RD_INTERVAL`      | `DefaultInterval` | 下载任务间隔时间（如 `30s`、`5m`）。                                                      |
| `-now`           | `RD_NOW`           | `false`           | 立即执行，不等待首次间隔。                                                                |
| `-once`          | `RD_ONCE`          | `false`           | 运行一次后退出。                                                                          |
| `-timeout`       | `RD_TIMEOUT`       | `DefaultTimeout`  | HTTP 客户端超时时间（如 `30s`、`2m`）。                                                   |
| `-asset_tag`     | `RD_ASSET_TAG`     | `false`           | 在文件名后追加版本标签（如 `file.zip` → `file-v1.0.0.zip`）。                             |
| `-asset_date`    | `RD_ASSET_DATE`    | `false`           | 在文件名后追加下载日期（如 `file.zip` → `file-20240502.zip`）。                           |
| `-asset_extract` | `RD_ASSET_EXTRACT` | `false`           | 自动解压下载的文件（支持 `.zip`、`.gz`、`.tar.gz`）。                                     |
| `-autoclean`     | `RD_AUTOCLEAN`     | `false`           | 下载新版本文件后，自动清理旧的 release 文件。                                             |

### 许可证

[Apache-2.0](LICENSE)
