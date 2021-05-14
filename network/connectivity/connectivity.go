package connectivity

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

var (
	ErrConnectivity = errors.New("failed connectivity check")
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Connectivity

// Connectivity has methods to check Internet connectivity.
type Connectivity interface {
	// Checks runs a DNS lookup, HTTP and HTTPs requests to all the domains given.
	// It returns any error encountered when doing so.
	Checks(ctx context.Context, domains ...string) (errs []error)
}

type connectivity struct {
	resolver *net.Resolver
	client   *http.Client
}

// NewConnectivity returns a new connectivity object.
func NewConnectivity(resolver *net.Resolver, client *http.Client) Connectivity {
	return &connectivity{
		resolver: resolver,
		client:   client,
	}
}

// Checks verifies the connection to the domains in terms of DNS, HTTP and HTTPS.
func (c *connectivity) Checks(ctx context.Context, domains ...string) (errs []error) {
	chErrors := make(chan []error)
	for _, domain := range domains {
		go func(domain string) {
			chErrors <- connectivityCheck(ctx, domain, c.resolver, c.client)
		}(domain)
	}
	for range domains {
		newErrors := <-chErrors
		errs = append(errs, newErrors...)
	}
	close(chErrors)
	return errs
}

func connectivityCheck(ctx context.Context, domain string,
	resolver *net.Resolver, client *http.Client) (errs []error) {
	chError := make(chan error)
	go func() { chError <- domainNameResolutionCheck(ctx, domain, resolver) }()
	go func() { chError <- httpGetCheck(ctx, "http://"+domain, client) }()
	go func() { chError <- httpGetCheck(ctx, "https://"+domain, client) }()
	for i := 0; i < 3; i++ {
		if err := <-chError; err != nil {
			err = fmt.Errorf("%w: for %s: %s", ErrConnectivity, domain, err)
			errs = append(errs, err)
		}
	}
	close(chError)
	return errs
}

var errNotOKHTTPStatus = errors.New("HTTP status is not OK")

func httpGetCheck(ctx context.Context, url string, client *http.Client) error {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %s", errNotOKHTTPStatus, response.Status)
	}
	return nil
}

func domainNameResolutionCheck(ctx context.Context, domain string, resolver *net.Resolver) error {
	_, err := resolver.LookupIP(ctx, "ip", domain)
	return err
}
