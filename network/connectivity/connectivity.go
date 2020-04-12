package connectivity

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/qdm12/golibs/network"
)

// Connectivity has methods to check Internet connectivity
//go:generate mockgen -destination=mock_connectivity/connectivity.go . Connectivity
type Connectivity interface {
	Checks(domains ...string) (errs []error)
}

type connectivity struct {
	checkDNS checkDNSFunc
	client   network.Client
}

// NewConnectivity returns a new connectivity object
func NewConnectivity(timeout time.Duration) Connectivity {
	return &connectivity{
		checkDNS: func(domain string) error {
			_, err := net.LookupIP(domain)
			return err
		},
		client: network.NewClient(timeout),
	}
}

type checkDNSFunc func(domain string) error

// Checks verifies the connection to the domains in terms of DNS, HTTP and HTTPS
func (c *connectivity) Checks(domains ...string) (errs []error) {
	chErrors := make(chan []error)
	for _, domain := range domains {
		go func(domain string) {
			chErrors <- connectivityCheck(domain, c.checkDNS, c.client)
		}(domain)
	}
	for range domains {
		newErrors := <-chErrors
		errs = append(errs, newErrors...)
	}
	close(chErrors)
	return errs
}

func connectivityCheck(domain string, checkDNS checkDNSFunc, client network.Client) (errs []error) {
	chError := make(chan error)
	go func() { chError <- domainNameResolutionCheck(domain, checkDNS) }()
	go func() { chError <- httpGetCheck("http://"+domain, client) }()
	go func() { chError <- httpGetCheck("https://"+domain, client) }()
	for i := 0; i < 3; i++ {
		if err := <-chError; err != nil {
			errs = append(errs, err)
		}
	}
	close(chError)
	return errs
}

func httpGetCheck(url string, client network.Client) error {
	_, status, err := client.GetContent(url)
	if err != nil {
		return fmt.Errorf("HTTP GET failed for %s: %w", url, err)
	} else if status != http.StatusOK {
		return fmt.Errorf("HTTP GET failed for %s: HTTP Status %d", url, status)
	}
	return nil
}

func domainNameResolutionCheck(domain string, checkDNS checkDNSFunc) error {
	if err := checkDNS(domain); err != nil {
		return fmt.Errorf("Domain name resolution is not working for %s: %w", domain, err)
	}
	return nil
}
