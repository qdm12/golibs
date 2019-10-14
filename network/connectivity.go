package network

import (
	"fmt"
	"net"
	"net/http"
)

// ConnectivityChecks verifies the connection to the domains in terms of DNS, HTTP and HTTPS
func ConnectivityChecks(client *http.Client, domains []string) (errs []error) {
	chErrors := make(chan []error)
	for _, domain := range domains {
		go connectivityCheck(client, domain, chErrors)
	}
	N := len(domains)
	for N > 0 {
		select {
		case errs := <-chErrors:
			errs = append(errs, errs...)
			N--
		}
	}
	close(chErrors)
	return errs
}

func connectivityCheck(client *http.Client, domain string, chErrors chan []error) {
	var errs []error
	chError := make(chan error)
	go domainNameResolutionCheckAsync(domain, chError)
	go httpGetCheckAsync(client, "http://"+domain, chError)
	go httpGetCheckAsync(client, "https://"+domain, chError)
	N := 3
	for N > 0 {
		select {
		case err := <-chError:
			if err != nil {
				errs = append(errs, err)
			}
			N--
		}
	}
	close(chError)
	chErrors <- errs
}

func httpGetCheckAsync(client *http.Client, URL string, chError chan error) {
	err := httpGetCheck(client, URL)
	chError <- err
}

func domainNameResolutionCheckAsync(domain string, chError chan error) {
	chError <- domainNameResolutionCheck(domain)
}

func httpGetCheck(client *http.Client, URL string) error {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return fmt.Errorf("HTTP GET failed for %s: %w", URL, err)
	}
	statusCode, _, err := DoHTTPRequest(client, req)
	if err != nil {
		return fmt.Errorf("HTTP GET failed for %s: %w", URL, err)
	} else if statusCode != 200 {
		return fmt.Errorf("HTTP GET failed for %s: HTTP Status %d", URL, statusCode)
	}
	return nil
}

func domainNameResolutionCheck(domain string) error {
	_, err := net.LookupIP(domain)
	if err != nil {
		return fmt.Errorf("Domain name resolution is not working for %s: %w", domain, err)
	}
	return nil
}