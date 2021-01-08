package clientip

import (
	"net"
	"net/http"
	"strings"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Extractor

type Extractor interface {
	HTTPRequest(r *http.Request) net.IP
}

type extractor struct {
	privateIPNets [8]net.IPNet
}

//nolint:gomnd
func NewExtractor() Extractor {
	return &extractor{
		privateIPNets: [8]net.IPNet{
			{ // localhost
				IP:   net.IP{127, 0, 0, 0},
				Mask: net.IPv4Mask(255, 0, 0, 0),
			},
			{ // 24-bit block
				IP:   net.IP{10, 0, 0, 0},
				Mask: net.IPv4Mask(255, 0, 0, 0),
			},
			{ // 20-bit block
				IP:   net.IP{172, 16, 0, 0},
				Mask: net.IPv4Mask(255, 240, 0, 0),
			},
			{ // 16-bit block
				IP:   net.IP{192, 168, 0, 0},
				Mask: net.IPv4Mask(255, 255, 0, 0),
			},
			{ // link local address
				IP:   net.IP{169, 254, 0, 0},
				Mask: net.IPv4Mask(255, 255, 0, 0),
			},
			{ // localhost IPv6
				IP:   net.IP{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1},
				Mask: net.IPMask{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
			},
			{ // unique local address IPv6
				IP:   net.IP{0xfc, 0x00, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{254, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
			{ // link local address IPv6
				IP:   net.IP{0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
				Mask: net.IPMask{255, 192, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
			},
		},
	}
}

func (e *extractor) HTTPRequest(r *http.Request) net.IP {
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
	publicXForwardedIPs := e.extractPublicIPs(xForwardedIPs)
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
