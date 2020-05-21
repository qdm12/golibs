ARG ALPINE_VERSION=3.11
ARG GO_VERSION=1.14

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder
ARG GOLANGCI_LINT_VERSION=v1.27.0
RUN apk --update add git
ENV CGO_ENABLED=0
RUN wget -O- -nv https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s ${GOLANGCI_LINT_VERSION}
WORKDIR /tmp/gobuild
COPY .golangci.yml .
COPY go.mod go.sum ./
RUN go mod download 2>&1
COPY . .
RUN go test ./...
RUN golangci-lint run
