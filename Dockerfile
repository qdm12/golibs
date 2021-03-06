ARG ALPINE_VERSION=3.12
ARG GO_VERSION=1.15

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
RUN apk --update add git
ENV CGO_ENABLED=0
WORKDIR /tmp/gobuild
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM --platform=$BUILDPLATFORM base AS test
# Note on the go race detector:
# - we set CGO_ENABLED=1 to have it enabled
# - we install g++ to support the race detector
ENV CGO_ENABLED=1
RUN apk --update --no-cache add g++

FROM --platform=$BUILDPLATFORM base AS lint
ARG GOLANGCI_LINT_VERSION=v1.36.0
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | \
  sh -s -- -b /usr/local/bin ${GOLANGCI_LINT_VERSION}
COPY .golangci.yml ./
RUN golangci-lint run --timeout=10m

FROM --platform=$BUILDPLATFORM base AS tidy
RUN git init && \
  git config user.email ci@localhost && \
  git config user.name ci && \
  git add -A && git commit -m ci && \
  sed -i '/\/\/ indirect/d' go.mod && \
  go mod tidy && \
  git diff --exit-code -- go.mod
