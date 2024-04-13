# **Experimental** Go libraries

When a Go library has a stable API and is useful across multiple projects, I extract it from this repository as a standalone repository.

<img height="250" src="https://raw.githubusercontent.com/qdm12/golibs/master/title.svg">

[![Build status](https://github.com/qdm12/golibs/workflows/CI/badge.svg?branch=master)](https://github.com/qdm12/golibs/actions?query=workflow%3A"CI")
[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/golibs.svg)](https://github.com/qdm12/golibs/commits/master)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/golibs.svg)](https://github.com/qdm12/golibs/graphs/contributors)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/golibs.svg)](https://github.com/qdm12/golibs/issues)

## Setup

```sh
go get github.com/qdm12/golibs
```

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
    1. You might want to use an editor such as [Visual Studio Code](https://code.visualstudio.com/download) with the [Go extension](https://code.visualstudio.com/docs/languages/go). Working settings are already in [.vscode/settings.json](https://github.com/qdm12/golibs/master/.vscode/settings.json).

    </p></details>

1. Commands available:

    ```sh
    # Lint the code
    golangci-lint run
    # Test the code
    go test ./...
    # Regenerate mocks for tests
    go generate -run "mockgen" ./...
    # Tidy modules dependencies
    go mod tidy
    # Run the CI steps with different Docker build targets:
    docker build --target lint .
    docker build --target test .
    docker build --target tidy .
    ```

1. See [Contributing](https://github.com/qdm12/golibs/master/.github/CONTRIBUTING.md) for more information on how to contribute to this repository.

## License

This repository is under an [MIT license](https://github.com/qdm12/golibs/master/license)
