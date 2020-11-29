# Reactserv

**WORK IN PROGRESS**

*HTTP server to serve React static files from memory. You can also set the root URL.*

<img height="200" src="https://raw.githubusercontent.com/qdm12/reactserv/master/title.svg?sanitize=true">

[![Build status](https://github.com/qdm12/reactserv/workflows/Buildx%20latest/badge.svg)](https://github.com/qdm12/reactserv/actions?query=workflow%3A%22Buildx+latest%22)
[![Docker Pulls](https://img.shields.io/docker/pulls/qmcgaw/reactserv.svg)](https://hub.docker.com/r/qmcgaw/reactserv)
[![Docker Stars](https://img.shields.io/docker/stars/qmcgaw/reactserv.svg)](https://hub.docker.com/r/qmcgaw/reactserv)
[![Image size](https://images.microbadger.com/badges/image/qmcgaw/reactserv.svg)](https://microbadger.com/images/qmcgaw/reactserv)
[![Image version](https://images.microbadger.com/badges/version/qmcgaw/reactserv.svg)](https://microbadger.com/images/qmcgaw/reactserv)

[![Join Slack channel](https://img.shields.io/badge/slack-@qdm12-yellow.svg?logo=slack)](https://join.slack.com/t/qdm12/shared_invite/enQtOTE0NjcxNTM1ODc5LTYyZmVlOTM3MGI4ZWU0YmJkMjUxNmQ4ODQ2OTAwYzMxMTlhY2Q1MWQyOWUyNjc2ODliNjFjMDUxNWNmNzk5MDk)
[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/reactserv.svg)](https://github.com/qdm12/reactserv/commits/master)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/reactserv.svg)](https://github.com/qdm12/reactserv/graphs/contributors)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/reactserv.svg)](https://github.com/qdm12/reactserv/issues)

## Features

- Reads the static React files from disk and serves them from memory
- Modify in-memory files with the `ROOT_URL` set, so using a reverse proxy is easier
- Compatible with `amd64`, `386`, `arm64`, `arm32v7`, `arm32v6`, `ppc64le` and `s390x` CPU architectures
- [Docker image tags and sizes](https://hub.docker.com/r/qmcgaw/reactserv/tags)

## Setup

1. Place your **compiled** React code in a directory, for example `/yourpath/react`.
1. Use the following command:

    ```sh
    docker run -d -p 8000:8000/tcp -v /yourpath/react:/srv:ro qmcgaw/reactserv
    ```

    You can also use [docker-compose.yml](https://github.com/qdm12/reactserv/blob/master/docker-compose.yml) with:

    ```sh
    docker-compose up -d
    ```

1. You can update the image with `docker pull qmcgaw/reactserv:latest` or use one of [tags available](https://hub.docker.com/r/qmcgaw/reactserv/tags)

### Environment variables

| Environment variable | Default | Possible values | Description |
| --- | --- | --- | --- |
| `LOG_ENCODING` | `console` | `json`, `console` | Logging format |
| `LOG_LEVEL` | `info` | `debug`, `info`, `warning`, `error` | Logging level |
| `LISTENING_PORT` | `8000` | Integer between `1` and `65535` | Internal listening TCP port |
| `ROOT_URL` | `/` | URL path *string* | URL path, used if behind a reverse proxy |
| `ROOT_DIR` | `srv` | Absolute or relative file path | Root filesystem path to read files from |
| `TZ` | `America/Montreal` | *string* | Timezone |

## Development

1. Setup your environment

    <details><summary>Using VSCode and Docker (easier)</summary><p>

    1. Install [Docker](https://docs.docker.com/install/)
       - On Windows, share a drive with Docker Desktop and have the project on that partition
       - On OSX, share your project directory with Docker Desktop
    1. With [Visual Studio Code](https://code.visualstudio.com/download), install the [remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
    1. In Visual Studio Code, press on `F1` and select `Remote-Containers: Open Folder in Container...`
    1. Your dev environment is ready to go!... and it's running in a container :+1: So you can discard it and update it easily!

    </p></details>

    <details><summary>Locally</summary><p>

    1. Install [Go](https://golang.org/dl/), [Docker](https://www.docker.com/products/docker-desktop) and [Git](https://git-scm.com/downloads)
    1. Install Go dependencies with

        ```sh
        go mod download
        ```

    1. Install [golangci-lint](https://github.com/golangci/golangci-lint#install)
    1. You might want to use an editor such as [Visual Studio Code](https://code.visualstudio.com/download) with the [Go extension](https://code.visualstudio.com/docs/languages/go). Working settings are already in [.vscode/settings.json](https://github.com/qdm12/reactserv/master/.vscode/settings.json).

    </p></details>

1. Commands available:

    ```sh
    # Build the binary
    go build cmd/app/main.go
    # Test the code
    go test ./...
    # Lint the code
    golangci-lint run
    # Build the Docker image
    docker build -t qmcgaw/reactserv .
    ```

1. See [Contributing](https://github.com/qdm12/reactserv/master/.github/CONTRIBUTING.md) for more information on how to contribute to this repository.

## TODOs

- Fix operation for index.html and directories
- Way to reload files into memory, maybe periodically?

## License

This repository is under an [MIT license](https://github.com/qdm12/reactserv/master/license) unless otherwise indicated
