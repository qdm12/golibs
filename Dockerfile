ARG ALPINE_VERSION=3.19
ARG GO_VERSION=1.22
ARG GOLANGCI_LINT_VERSION=v1.56.2

FROM --platform=${BUILDPLATFORM} qmcgaw/binpot:golangci-lint-${GOLANGCI_LINT_VERSION} AS golangci-lint

FROM --platform=${BUILDPLATFORM} golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS base
RUN apk --update --no-cache add git g++ findutils
ENV CGO_ENABLED=0
WORKDIR /tmp/gobuild
COPY --from=golangci-lint /bin /go/bin/golangci-lint
COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM --platform=${BUILDPLATFORM} base AS test
# Note on the go race detector:
# - we set CGO_ENABLED=1 to have it enabled
# - we installed g++ to support the race detector
ENV CGO_ENABLED=1
RUN touch coverage.txt
ENTRYPOINT go test -race -coverprofile=coverage.txt -covermode=atomic ./...

FROM --platform=${BUILDPLATFORM} base AS lint
COPY .golangci.yml ./
RUN golangci-lint run --timeout=10m

FROM --platform=${BUILDPLATFORM} base AS uptodate
RUN git init && \
  git config user.email ci@localhost && \
  git config user.name ci && \
  git config core.fileMode false && \
  git add -A && \
  git commit -m snapshot && \
  # Check modules are tidied
  go mod tidy && \
  git diff --exit-code -- go.mod && \
  # Mocks
  grep -lr -E '^// Code generated by MockGen\. DO NOT EDIT\.$' . | xargs -r -d '\n' rm && \
  go generate -run "mockgen" ./... && \
  git diff --exit-code && \
  rm -rf .git/
