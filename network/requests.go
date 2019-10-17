package network

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/qdm12/golibs/security"
)

var userAgents = []string{
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Safari/537.36",
	"Mozilla/5.0 (Linux; Android 8.0.0; Nexus 5X Build/OPR4.170623.006) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Mobile Safari/537.36",
	"Mozilla/5.0 (iPad; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.34 (KHTML, like Gecko) Version/11.0 Mobile/15A5341f Safari/604.1",
	"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1",
}

// GetContentParamsType contains the optional parameters to use
// to get content with an HTTP get
type GetContentParamsType struct {
	DisguisedUserAgent bool
}

// DoHTTPRequest performs an HTTP request and returns the status, content and eventual error
func DoHTTPRequest(client *http.Client, request *http.Request) (status int, content []byte, err error) {
	response, err := client.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return status, nil, err
	}
	content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return status, nil, err
	}
	return response.StatusCode, content, nil
}

// GetContent returns the content and eventual error from an HTTP GET to a given URL
func GetContent(httpClient *http.Client, URL string, params GetContentParamsType) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot GET content of URL %s: %w", URL, err)
	}
	if params.DisguisedUserAgent {
		req.Header.Set("User-Agent", userAgents[security.GenerateRandomInt(len(userAgents))])
	}
	status, content, err := DoHTTPRequest(httpClient, req)
	if err != nil {
		return nil, fmt.Errorf("cannot GET content of URL %s: %w", URL, err)
	}
	if status != 200 {
		return nil, fmt.Errorf("cannot GET content of URL %s (status %d)", URL, status)
	}
	return content, nil
}
