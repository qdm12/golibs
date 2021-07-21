# Golibs

*Golang libraries I use across my projects*

<img height="250" src="https://raw.githubusercontent.com/qdm12/golibs/master/title.svg?sanitize=true">

[![Build status](https://github.com/qdm12/golibs/workflows/CI/badge.svg?branch=master)](https://github.com/qdm12/golibs/actions?query=workflow%3A"CI")
[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/golibs.svg)](https://github.com/qdm12/golibs/commits/master)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/golibs.svg)](https://github.com/qdm12/golibs/graphs/contributors)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/golibs.svg)](https://github.com/qdm12/golibs/issues)

## Setup

Simply import one of the following libraries in your Go code:

- `"github.com/qdm12/golibs/admin"` for supervising your application, it contains a Gotify client
- `"github.com/qdm12/golibs/command"` for interacting the shell command line
- `"github.com/qdm12/golibs/crypto"` for encryption, randomness and checksum functions
- `"github.com/qdm12/golibs/files"` to interact with the filesystem
- `"github.com/qdm12/golibs/format"` to format things to strings
- `"github.com/qdm12/golibs/logging"` for logging functions with a global Zap logger
- `"github.com/qdm12/golibs/network"` for HTTP requests, IP address processing and connectivity checks
- `"github.com/qdm12/golibs/os"` for OS related operations like file manipulation.
- `"github.com/qdm12/golibs/params"` for parsing and verifying parameters from environment variables
- `"github.com/qdm12/golibs/redis"` for basic initialization of a Redis database pool of connections
- `"github.com/qdm12/golibs/verification"` for verification functions such as email checking or regex based checking.

For each package, some mocks are generated using [mockgen](https://github.com/golang/mock#running-mockgen) and can be imported with, for example

```go
import github.com/qdm12/golibs/verification/mock_verification
```

and used with [gomock](https://github.com/golang/mock#building-mocks) for testing.

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
    # Build the binary
    go build cmd/app/main.go
    # Test the code
    go test ./...
    # Regenerate mocks for tests
    go generate ./...
    # Lint the code
    golangci-lint run
    # Build the Docker image to run tests and linting
    docker build .
    ```

1. See [Contributing](https://github.com/qdm12/golibs/master/.github/CONTRIBUTING.md) for more information on how to contribute to this repository.

## TODOs

- HTTP server/client unit tests
- Server rework to write unique request ID (see timesheet)
- More hashing functions

## License

This repository is under an [MIT license](https://github.com/qdm12/golibs/master/license)
