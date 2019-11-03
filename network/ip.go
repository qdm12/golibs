package network

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/qdm12/golibs/verification"
	"go.uber.org/zap"
)

var privateCIDRs []*net.IPNet

func init() {
	privateCIDRBlocks := [8]string{
		"127.0.0.1/8",    // localhost
		"10.0.0.0/8",     // 24-bit block
		"172.16.0.0/12",  // 20-bit block
		"192.168.0.0/16", // 16-bit block
		"169.254.0.0/16", // link local address
		"::1/128",        // localhost IPv6
		"fc00::/7",       // unique local address IPv6
		"fe80::/10",      // link local address IPv6
	}
	for _, privateCIDRBlock := range privateCIDRBlocks {
		_, CIDR, _ := net.ParseCIDR(privateCIDRBlock)
		privateCIDRs = append(privateCIDRs, CIDR)
	}
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
func GetClientIPHeaders(r *http.Request) (headers IPHeaders) {
	if r == nil {
		return headers
	}
	headers.RemoteAddress = strings.ReplaceAll(r.RemoteAddr, " ", "")
	headers.XRealIP = strings.ReplaceAll(r.Header.Get("X-Real-Ip"), " ", "")
	headers.XForwardedFor = strings.ReplaceAll(r.Header.Get("X-Forwarded-For"), " ", "")
	return headers
}

// GetClientIP returns one single client IP address
func GetClientIP(r *http.Request) (IP string, err error) {
	headers := GetClientIPHeaders(r)
	if headers.isVoid() {
		return "", fmt.Errorf("no IP address found in client request")
	}
	// Extract relevant IP data from headers
	remoteIP, err := getRemoteIP(headers.RemoteAddress)
	if err != nil {
		return "", err
	}
	// No headers so it can only be RemoteAddress
	if headers.XRealIP == "" && headers.XForwardedFor == "" {
		return remoteIP, nil
	}
	// 3. RemoteAddress is the proxy server forwarding the IP so
	// we look into the HTTP headers to get the client IP
	xForwardedIPs := getXForwardedIPs(headers.XForwardedFor)
	// TODO check number of ips to match number of proxies setup
	publicXForwardedIPs := extractPublicIPs(xForwardedIPs)
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

func ipIsValid(IP string) bool {
	netIP := net.ParseIP(IP)
	return netIP != nil
}

func netIPIsPrivate(netIP net.IP) bool {
	for i := range privateCIDRs {
		if privateCIDRs[i].Contains(netIP) {
			return true
		}
	}
	return false
}

func splitHostPort(address string) (IP, port string, err error) {
	if strings.ContainsRune(address, '[') && strings.ContainsRune(address, ']') {
		// should be an IPv6 address with brackets
		return net.SplitHostPort(address)
	}
	if strings.Count(address, ":") > 1 { // could be an IPv6 without brackets
		i := strings.LastIndex(address, ":")
		port = address[i+1 : len(address)]
		IP = address[0:i]
		if !ipIsValid(IP) {
			return net.SplitHostPort(address)
		}
		return IP, port, nil
	}
	// IPv4 address
	return net.SplitHostPort(address)
}

func getRemoteIP(remoteAddr string) (IP string, err error) {
	IP = remoteAddr
	if strings.ContainsRune(IP, ':') {
		var port string
		IP, port, err = splitHostPort(IP)
		if err != nil {
			return "", err
		}
		if len(port) > 0 {
			if err := verification.VerifyPort(port); err != nil {
				return "", fmt.Errorf("remote address %q is invalid: %w", remoteAddr, err)
			}
		}
	}
	if !ipIsValid(IP) {
		return "", fmt.Errorf("IP address %q is not valid", IP)
	}
	return IP, nil
}

func extractPublicIPs(IPs []string) (publicIPs []string) {
	for _, IP := range IPs {
		if !ipIsValid(IP) {
			zap.L().Warn("IP address is not valid", zap.String("IP", IP))
			continue
		}
		netIP := net.ParseIP(IP)
		private := netIPIsPrivate(netIP)
		if !private {
			publicIPs = append(publicIPs, IP)
		}
	}
	return publicIPs
}

func getXForwardedIPs(XForwardedFor string) (IPs []string) {
	if len(XForwardedFor) == 0 {
		return nil
	}
	XForwardedFor = strings.ReplaceAll(XForwardedFor, " ", "")
	XForwardedFor = strings.ReplaceAll(XForwardedFor, "\t", "")
	XForwardForIPs := strings.Split(XForwardedFor, ",")
	for _, IP := range XForwardForIPs {
		if !ipIsValid(IP) {
			zap.L().Warn("IP address is not valid", zap.String("IP", IP))
			continue
		}
		IPs = append(IPs, IP)
	}
	return IPs
}
