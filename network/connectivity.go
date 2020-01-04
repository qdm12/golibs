package network

import (
	"fmt"
	"net"
	"time"
)

// Connectivity has methods to check Internet connectivity
type Connectivity interface {
	Checks(domains ...string) (errs []error)
}

// ConnectivityImpl implements Connectivity
type ConnectivityImpl struct {
	checkDNS checkDNSFunc
	client   Client
}

// NewConnectivity returns a new connectivity object
func NewConnectivity(timeout time.Duration) Connectivity {
	return &ConnectivityImpl{
		checkDNS: func(domain string) error {
			_, err := net.LookupIP(domain)
			return err
		},
		client: NewClient(timeout),
	}
}

type checkDNSFunc func(domain string) error

// Checks verifies the connection to the domains in terms of DNS, HTTP and HTTPS
func (c *ConnectivityImpl) Checks(domains ...string) (errs []error) {
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

func connectivityCheck(domain string, checkDNS checkDNSFunc, client Client) (errs []error) {
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

func httpGetCheck(URL string, client Client) error {
	_, status, err := client.GetContent(URL)
	if err != nil {
		return fmt.Errorf("HTTP GET failed for %s: %w", URL, err)
	} else if status != 200 {
		return fmt.Errorf("HTTP GET failed for %s: HTTP Status %d", URL, status)
	}
	return nil
}

func domainNameResolutionCheck(domain string, checkDNS checkDNSFunc) error {
	if err := checkDNS(domain); err != nil {
		return fmt.Errorf("Domain name resolution is not working for %s: %w", domain, err)
	}
	return nil
}
