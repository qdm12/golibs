package connectivity

import (
	"context"
	"net/http"
	urlpkg "net/url"
)

// NewHTTPSGetChecker creates a new HTTPs checker which expects
// the HTTP status given when doing an HTTP GET request.
func NewHTTPSGetChecker(client *http.Client, expectedStatus int) *HTTPSGetChecker {
	return &HTTPSGetChecker{
		client:         client,
		expectedStatus: expectedStatus,
	}
}

// HTTPSGetChecker implements a checker to send HTTPs GET requests
// and verify the response status code against the one specified.
type HTTPSGetChecker struct {
	client         *http.Client
	expectedStatus int
}

// ParallelChecks verifies the connectivity to each of the urls
// using an HTTPs GET request and comparing the received status code.
// It returns a slice of errors with the same indexing and order as the
// urls, meaning that some errors might be nil or not. You should ensure
// to iterate over the errors and check each of them.
func (c *HTTPSGetChecker) ParallelChecks(ctx context.Context, urls []string) (errs []error) {
	return parallelChecks(ctx, c, urls)
}

// Check verifies the HTTP response status code matches the expected
// HTTP status code when sending an HTTP GET request to the url over HTTPS.
func (c *HTTPSGetChecker) Check(ctx context.Context, url string) error {
	u, err := urlpkg.Parse(url)
	if err != nil {
		return err
	}

	u.Scheme = "https"
	return httpGetCheck(ctx, c.client, u.String(), c.expectedStatus)
}
