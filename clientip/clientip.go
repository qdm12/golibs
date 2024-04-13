package clientip

import (
	"net"
	"net/http"
	"strings"
)

type Parser struct {
	privateIPNets [8]net.IPNet
}

func NewParser() *Parser {
	return &Parser{
		privateIPNets: privateIPNets(),
	}
}

func (p *Parser) ParseHTTPRequest(r *http.Request) net.IP {
	if r == nil {
		return nil
	}

	remoteAddress := removeSpaces(r.RemoteAddr)
	xRealIP := removeSpaces(r.Header.Get("X-Real-IP"))
	xForwardedFor := r.Header.Values("X-Forwarded-For")
	for i := range xForwardedFor {
		xForwardedFor[i] = removeSpaces(xForwardedFor[i])
	}

	// No header so it can only be remoteAddress
	if xRealIP == "" && len(xForwardedFor) == 0 {
		return getIPFromHostPort(remoteAddress)
	}

	// remoteAddress is the last proxy server forwarding the traffic
	// so we look into the HTTP headers to get the client IP
	xForwardedIPs := parseIPs(xForwardedFor)
	publicXForwardedIPs := p.extractPublicIPs(xForwardedIPs)
	if len(publicXForwardedIPs) > 0 {
		// first public XForwardedIP should be the client IP
		return publicXForwardedIPs[0]
	}

	// If all forwarded IP addresses are private we use the x-real-ip
	// address if it exists
	if xRealIP != "" {
		return getIPFromHostPort(xRealIP)
	}

	// Client IP is the first private IP address in the chain
	return xForwardedIPs[0]
}

func removeSpaces(header string) string {
	header = strings.ReplaceAll(header, " ", "")
	header = strings.ReplaceAll(header, "\t", "")
	return header
}
