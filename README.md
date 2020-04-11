# Golibs

*Golang libraries I use across my projects*

[![golibs](https://github.com/qdm12/golibs/raw/master/title.png)](https://hub.docker.com/r/qmcgaw/REPONAME_DOCKER)

[![Build Status](https://travis-ci.org/qdm12/golibs.svg?branch=master)](https://travis-ci.org/qdm12/golibs)

[![GitHub last commit](https://img.shields.io/github/last-commit/qdm12/golibs.svg)](https://github.com/qdm12/golibs/issues)
[![GitHub commit activity](https://img.shields.io/github/commit-activity/y/qdm12/golibs.svg)](https://github.com/qdm12/golibs/issues)
[![GitHub issues](https://img.shields.io/github/issues/qdm12/golibs.svg)](https://github.com/qdm12/golibs/issues)

## Setup

Simply import one of the following libraries in your Go code:

- `"github.com/qdm12/golibs/admin"` for supervising your application, it contains a Gotify client
- `"github.com/qdm12/golibs/database"` for basic initialization of a postgreSQL database pool of connections
- `"github.com/qdm12/golibs/healthcheck"` for server and client functions to provide a healthcheck
- `"github.com/qdm12/golibs/logging"` for logging functions with a global Zap logger
- `"github.com/qdm12/golibs/network"` for HTTP requests, IP address processing and connectivity checks
- `"github.com/qdm12/golibs/params"` for parsing and verifying parameters from environment variables
- `"github.com/qdm12/golibs/redis"` for basic initialization of a Redis database pool of connections
- `"github.com/qdm12/golibs/security"` for encryption, randomness and checksum functions
- `"github.com/qdm12/golibs/server"` for HTTP server functions
- `"github.com/qdm12/golibs/signals"` for termination signal catching for a graceful shutdown of the application
- `"github.com/qdm12/golibs/verification"` for verification functions such as email checking or regex based checking.

## Development

### Using VSCode and Docker

1. Install [Docker](https://docs.docker.com/install/)
    - On Windows, share a drive with Docker Desktop and have the project on that partition
1. With [Visual Studio Code](https://code.visualstudio.com/download), install the [remote containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
1. In Visual Studio Code, press on `F1` and select `Remote-Containers: Open Folder in Container...`
1. Your dev environment is ready to go!... and it's running in a container :+1:

Regenerate mocks for tests using `go generate ./....`

## TODOs

- HTTP server/client unit tests
- Server rework to write unique request ID (see timesheet)
- More hashing functions

## License

This repository is under an [MIT license](https://github.com/qdm12/golibs/master/license)
