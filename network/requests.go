package network

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/qdm12/golibs/security"
)

// Client has methods to do HTTP requests as a client
type Client interface {
	DoHTTPRequest(request *http.Request) (status int, content []byte, err error)
	GetContent(URL string, setters ...GetContentSetter) (content []byte, status int, err error)
	Close()
}

type httpClient interface {
	Do(r *http.Request) (*http.Response, error)
	CloseIdleConnections()
}

func newHTTPClient(timeout time.Duration) httpClient {
	return &http.Client{Timeout: timeout}
}

// ClientImpl is the implementation for IClient
type ClientImpl struct {
	httpClient httpClient
	readBody   func(r io.Reader) ([]byte, error)
	userAgents []string
	random     security.Random
}

// NewClient creates a new easy to use HTTP client
func NewClient(timeout time.Duration) Client {
	return &ClientImpl{
		httpClient: newHTTPClient(timeout),
		readBody:   ioutil.ReadAll,
		userAgents: []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Safari/537.36",
			"Mozilla/5.0 (Linux; Android 8.0.0; Nexus 5X Build/OPR4.170623.006) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Mobile Safari/537.36",
			"Mozilla/5.0 (iPad; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.34 (KHTML, like Gecko) Version/11.0 Mobile/15A5341f Safari/604.1",
			"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1",
		},
	}
}

// Close terminates idle connections of the HTTP client
func (c *ClientImpl) Close() {
	if c.httpClient != nil {
		c.httpClient.CloseIdleConnections()
	}
}

// DoHTTPRequest performs an HTTP request and returns the status, content and eventual error
func (c *ClientImpl) DoHTTPRequest(request *http.Request) (status int, content []byte, err error) {
	response, err := c.httpClient.Do(request)
	if response != nil {
		defer response.Body.Close()
	}
	if err != nil {
		return status, nil, err
	}
	content, err = c.readBody(response.Body)
	if err != nil {
		return status, nil, err
	}
	return response.StatusCode, content, nil
}

type getContentOptions struct {
	randomUserAgent bool
}

// GetContentSetter is a setter for options to GetContent
type GetContentSetter func(options *getContentOptions)

// UseRandomUserAgent sets a random realistic user agent to the GetContent HTTP request
func UseRandomUserAgent() GetContentSetter {
	return func(options *getContentOptions) {
		options.randomUserAgent = true
	}
}

// GetContent returns the content and eventual error from an HTTP GET to a given URL
func (c *ClientImpl) GetContent(URL string, setters ...GetContentSetter) (content []byte, status int, err error) {
	var options getContentOptions
	for _, setter := range setters {
		setter(&options)
	}
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, status, fmt.Errorf("cannot GET content of URL %s: %w", URL, err)
	}
	if options.randomUserAgent {
		req.Header.Set("User-Agent", c.userAgents[c.random.GenerateRandomInt(len(c.userAgents))])
	}
	status, content, err = c.DoHTTPRequest(req)
	if err != nil {
		return nil, status, fmt.Errorf("cannot GET content of URL %s: %w", URL, err)
	}
	return content, status, nil
}
