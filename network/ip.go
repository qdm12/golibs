package network

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
)

type IPManager interface {
	GetClientIPHeaders(r *http.Request) (headers IPHeaders)
	GetClientIP(r *http.Request) (ip string, err error)
}

type ipManager struct {
	logger   logging.Logger
	verifier verification.Verifier
}

func NewIPManager(logger logging.Logger) IPManager {
	return &ipManager{
		logger, verification.NewVerifier()}
}

// IPHeaders contains all the raw IP headers of an HTTP request
type IPHeaders struct {
	RemoteAddress string
	XRealIP       string
	XForwardedFor string
}

func (headers *IPHeaders) String() string {
	if headers == nil {
		return "remoteAddr= | xRealIP= | xForwardedFor="
	}
	return fmt.Sprintf("remoteAddr=%q | xRealIP=%q | xForwardedFor=%q",
		headers.RemoteAddress, headers.XRealIP, headers.XForwardedFor)
}

func (headers *IPHeaders) isVoid() bool {
	return headers == nil || (headers.RemoteAddress == "" &&
		headers.XRealIP == "" &&
		headers.XForwardedFor == "")
}

// GetClientIPHeaders returns the IP related HTTP headers from a request
func (m *ipManager) GetClientIPHeaders(r *http.Request) (headers IPHeaders) {
	if r == nil {
		return headers
	}
	headers.RemoteAddress = strings.ReplaceAll(r.RemoteAddr, " ", "")
	headers.XRealIP = strings.ReplaceAll(r.Header.Get("X-Real-Ip"), " ", "")
	headers.XForwardedFor = strings.ReplaceAll(r.Header.Get("X-Forwarded-For"), " ", "")
	return headers
}

// GetClientIP returns one single client IP address
func (m *ipManager) GetClientIP(r *http.Request) (ip string, err error) {
	headers := m.GetClientIPHeaders(r)
	if headers.isVoid() {
		return "", fmt.Errorf("no IP address found in client request")
	}
	// Extract relevant IP data from headers
	remoteIP, err := getRemoteIP(m.verifier.VerifyPort, headers.RemoteAddress)
	if err != nil {
		return "", err
	}
	// No headers so it can only be RemoteAddress
	if headers.XRealIP == "" && headers.XForwardedFor == "" {
		return remoteIP, nil
	}
	// 3. RemoteAddress is the proxy server forwarding the IP so
	// we look into the HTTP headers to get the client IP
	xForwardedIPs, warnings := getXForwardedIPs(headers.XForwardedFor)
	for _, warning := range warnings {
		m.logger.Warn(warning)
	}
	// TODO check number of ips to match number of proxies setup
	publicXForwardedIPs, warnings := extractPublicIPs(xForwardedIPs)
	for _, warning := range warnings {
		m.logger.Warn(warning)
	}
	if len(publicXForwardedIPs) > 0 {
		// first XForwardedIP should be the client IP
		return publicXForwardedIPs[0], nil
	}
	if headers.XRealIP != "" {
		if ipIsValid(headers.XRealIP) {
			return headers.XRealIP, nil
		}
	}
	// latest private XForwardedFor IP
	if len(xForwardedIPs) > 0 {
		return xForwardedIPs[len(xForwardedIPs)-1], nil
	}
	return remoteIP, nil
}

func ipIsValid(ip string) bool {
	netIP := net.ParseIP(ip)
	return netIP != nil
}

func netIPIsPrivate(netIP net.IP) bool {
	for _, privateCIDRBlock := range [8]string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	} {
		_, CIDR, _ := net.ParseCIDR(privateCIDRBlock)
		if CIDR.Contains(netIP) {
			return true
		}
	}
	return false
}

func splitHostPort(address string) (ip, port string, err error) {
	if strings.ContainsRune(address, '[') && strings.ContainsRune(address, ']') {
		// should be an IPv6 address with brackets
		return net.SplitHostPort(address)
	}
	const ipv4MaxColons = 1
	if strings.Count(address, ":") > ipv4MaxColons { // could be an IPv6 without brackets
		i := strings.LastIndex(address, ":")
		port = address[i+1:]
		ip = address[0:i]
		if !ipIsValid(ip) {
			return net.SplitHostPort(address)
		}
		return ip, port, nil
	}
	// IPv4 address
	return net.SplitHostPort(address)
}

func getRemoteIP(verifyPort func(port string) error, remoteAddr string) (ip string, err error) {
	ip = remoteAddr
	if strings.ContainsRune(ip, ':') {
		var port string
		ip, port, err = splitHostPort(ip)
		if err != nil {
			return "", err
		}
		if len(port) > 0 {
			if err := verifyPort(port); err != nil {
				return "", fmt.Errorf("remote address %q is invalid: %w", remoteAddr, err)
			}
		}
	}
	if !ipIsValid(ip) {
		return "", fmt.Errorf("IP address %q is not valid", ip)
	}
	return ip, nil
}

func extractPublicIPs(ips []string) (publicIPs []string, warnings []string) {
	for _, IP := range ips {
		if !ipIsValid(IP) {
			warnings = append(warnings, fmt.Sprintf("IP address %q is not valid", IP))
			continue
		}
		netIP := net.ParseIP(IP)
		private := netIPIsPrivate(netIP)
		if !private {
			publicIPs = append(publicIPs, IP)
		}
	}
	return publicIPs, warnings
}

func getXForwardedIPs(xForwardedFor string) (ips []string, warnings []string) {
	if len(xForwardedFor) == 0 {
		return nil, nil
	}
	xForwardedFor = strings.ReplaceAll(xForwardedFor, " ", "")
	xForwardedFor = strings.ReplaceAll(xForwardedFor, "\t", "")
	XForwardForIPs := strings.Split(xForwardedFor, ",")
	for _, IP := range XForwardForIPs {
		if !ipIsValid(IP) {
			warnings = append(warnings, fmt.Sprintf("IP address %q is not valid", IP))
			continue
		}
		ips = append(ips, IP)
	}
	return ips, warnings
}
