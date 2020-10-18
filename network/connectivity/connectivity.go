package connectivity

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/qdm12/golibs/network"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Connectivity

// Connectivity has methods to check Internet connectivity.
type Connectivity interface {
	// Checks runs a DNS lookup, HTTP and HTTPs requests to all the domains given.
	// It returns any error encountered when doing so.
	Checks(ctx context.Context, domains ...string) (errs []error)
}

type connectivity struct {
	checkDNS checkDNSFunc
	client   network.Client
}

// NewConnectivity returns a new connectivity object.
func NewConnectivity(timeout time.Duration) Connectivity {
	return &connectivity{
		checkDNS: func(ctx context.Context, host string) error {
			_, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
			return err
		},
		client: network.NewClient(timeout),
	}
}

type checkDNSFunc func(ctx context.Context, host string) error

// Checks verifies the connection to the domains in terms of DNS, HTTP and HTTPS.
func (c *connectivity) Checks(ctx context.Context, domains ...string) (errs []error) {
	chErrors := make(chan []error)
	for _, domain := range domains {
		go func(domain string) {
			chErrors <- connectivityCheck(ctx, domain, c.checkDNS, c.client)
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
	checkDNS checkDNSFunc, client network.Client) (errs []error) {
	chError := make(chan error)
	go func() { chError <- domainNameResolutionCheck(ctx, domain, checkDNS) }()
	go func() { chError <- httpGetCheck(ctx, "http://"+domain, client) }()
	go func() { chError <- httpGetCheck(ctx, "https://"+domain, client) }()
	for i := 0; i < 3; i++ {
		if err := <-chError; err != nil {
			errs = append(errs, err)
		}
	}
	close(chError)
	return errs
}

func httpGetCheck(ctx context.Context, url string, client network.Client) error {
	_, status, err := client.Get(ctx, url)
	if err != nil {
		return fmt.Errorf("HTTP GET failed for %s: %w", url, err)
	} else if status != http.StatusOK {
		return fmt.Errorf("HTTP GET failed for %s: HTTP Status %d", url, status)
	}
	return nil
}

func domainNameResolutionCheck(ctx context.Context, domain string, checkDNS checkDNSFunc) error {
	if err := checkDNS(ctx, domain); err != nil {
		return fmt.Errorf("Domain name resolution is not working for %s: %w", domain, err)
	}
	return nil
}
