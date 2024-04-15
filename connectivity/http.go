package connectivity

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	urlpkg "net/url"
)

// NewHTTPGetChecker creates a new HTTPs checker which expects
// the HTTP status given when doing an HTTP GET request.
func NewHTTPGetChecker(client *http.Client, expectedStatus int) *HTTPGetChecker {
	return &HTTPGetChecker{
		client:         client,
		expectedStatus: expectedStatus,
	}
}

// HTTPGetChecker implements a checker to send HTTP GET requests
// and verify the response status code against the one specified.
type HTTPGetChecker struct {
	client         *http.Client
	expectedStatus int
}

// ParallelChecks verifies the connectivity to each of the urls
// using a plaintext HTTP GET request and comparing the received status code.
// It returns a slice of errors with the same indexing and order as the
// urls, meaning that some errors might be nil or not. You should ensure
// to iterate over the errors and check each of them.
func (c *HTTPGetChecker) ParallelChecks(ctx context.Context, urls []string) (errs []error) {
	return parallelChecks(ctx, c, urls)
}

// Check verifies the HTTP response status code matches the expected
// HTTP status code when sending an HTTP GET request to the url over
// plaintext HTTP.
func (c *HTTPGetChecker) Check(ctx context.Context, url string) error {
	u, err := urlpkg.Parse(url)
	if err != nil {
		return fmt.Errorf("parsing url: %w", err)
	}

	u.Scheme = "http"
	return httpGetCheck(ctx, c.client, u.String(), c.expectedStatus)
}

var ErrHTTPStatusUnexpected = errors.New("unexpected HTTP status received")

func httpGetCheck(ctx context.Context, client *http.Client,
	url string, expectedStatus int) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("doing request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != expectedStatus {
		return fmt.Errorf("%w: expected %d and received %s",
			ErrHTTPStatusUnexpected, expectedStatus, response.Status)
	}

	return nil
}
