package network

import "net/http"

// Only used for tests

//go:generate mockgen -destination=mockHttpClient_test.go -package=network github.com/qdm12/golibs/network HTTPClient
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
	CloseIdleConnections()
}
