package connectivity

import (
	"context"
	"errors"
	"net"
	urlpkg "net/url"
)

var errDNSResolveNoIP = errors.New("DNS resolution resulted in no IP address")

// NewDNSChecker creates a new DNS checker.
func NewDNSChecker(resolver *net.Resolver) *DNSChecker {
	return &DNSChecker{
		resolver: resolver,
	}
}

// DNSChecker implements a checker to DNS resolve domain names to
// IP addresses and verify it has at least one IP address matching.
type DNSChecker struct {
	resolver *net.Resolver
}

// ParallelChecks verifies the domain name of each of the urls
// resolves successfully to at least one IP address.
// It returns a slice of errors with the same indexing and order as the
// urls, meaning that some errors might be nil or not. You should ensure
// to iterate over the errors and check each of them.
func (c *DNSChecker) ParallelChecks(ctx context.Context, urls []string) (errs []error) {
	return parallelChecks(ctx, c, urls)
}

// Check verifies the domain name of the urls resolves successfully
// to at least one IP address.
func (c *DNSChecker) Check(ctx context.Context, url string) error {
	u, err := urlpkg.Parse(url)
	if err != nil {
		return err
	}

	domain := u.Hostname()
	ips, err := c.resolver.LookupIP(ctx, "ip", domain)
	if err != nil {
		return err
	} else if len(ips) == 0 {
		return errDNSResolveNoIP
	}

	return err
}
