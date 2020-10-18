package network

import (
	"context"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/qdm12/golibs/crypto/random"
	"golang.org/x/net/context/ctxhttp"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Client

// Client has methods to do HTTP requests as a client
type Client interface {
	// Do runs the given request and returns an HTTP status,
	// the data content and an eventual error
	Do(ctx context.Context, request *http.Request) (content []byte, status int, err error)
	// Get runs an HTTP GET operation at a given URL and returns the content, status and error
	Get(ctx context.Context, URL string, setters ...GetSetter) (content []byte, status int, err error)
	// Close closes any idle connections remaining for this client
	Close()
}

type client struct {
	httpClient *http.Client
	userAgents []string
	random     random.Random
}

// NewClient creates a new easy to use HTTP client
func NewClient(timeout time.Duration) Client {
	return &client{
		httpClient: &http.Client{Timeout: timeout},
		userAgents: []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_5) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/50.0.2661.102 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Safari/537.36",
			"Mozilla/5.0 (Linux; Android 8.0.0; Nexus 5X Build/OPR4.170623.006) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/65.0.3325.146 Mobile Safari/537.36",
			"Mozilla/5.0 (iPad; CPU OS 11_0 like Mac OS X) AppleWebKit/604.1.34 (KHTML, like Gecko) Version/11.0 Mobile/15A5341f Safari/604.1",
			"Mozilla/5.0 (iPhone; CPU iPhone OS 11_0 like Mac OS X) AppleWebKit/604.1.38 (KHTML, like Gecko) Version/11.0 Mobile/15A372 Safari/604.1",
		},
		random: random.NewRandom(),
	}
}

// Close terminates idle connections of the HTTP client
func (c *client) Close() {
	c.httpClient.CloseIdleConnections()
}

// Do performs an HTTP request and returns the status, content and eventual error
func (c *client) Do(ctx context.Context, request *http.Request) (content []byte, status int, err error) {
	response, err := ctxhttp.Do(ctx, c.httpClient, request)
	if err != nil {
		return nil, 0, err
	}
	defer response.Body.Close()
	content, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, response.StatusCode, err
	}
	return content, response.StatusCode, nil
}

type getOptions struct {
	randomUserAgent bool
}

// GetContentSetter is a setter for options to GetContent
type GetSetter func(options *getOptions)

// UseRandomUserAgent sets a random realistic user agent to the GetContent HTTP request
func UseRandomUserAgent() GetSetter {
	return func(options *getOptions) {
		options.randomUserAgent = true
	}
}

// GetContent returns the content and eventual error from an HTTP GET to a given URL
func (c *client) Get(ctx context.Context, url string, setters ...GetSetter) (content []byte, status int, err error) {
	var options getOptions
	for _, setter := range setters {
		setter(&options)
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, 0, err
	}
	if options.randomUserAgent {
		req.Header.Set("User-Agent", c.userAgents[c.random.GenerateRandomInt(len(c.userAgents))])
	}
	return c.Do(ctx, req)
}
