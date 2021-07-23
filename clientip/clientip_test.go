package clientip

import (
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewParser(t *testing.T) {
	t.Parallel()
	expectedParser := &Parser{
		privateIPNets: privateIPNets(),
	}

	parser := NewParser()
	assert.Equal(t, expectedParser, parser)
}

func Test_Parser_ParseHTTPRequest(t *testing.T) {
	t.Parallel()

	makeHeader := func(keyValues map[string][]string) http.Header {
		header := http.Header{}
		for key, values := range keyValues {
			for _, value := range values {
				header.Add(key, value)
			}
		}
		return header
	}

	testCases := map[string]struct {
		r  *http.Request
		ip net.IP
	}{
		"nil request": {},
		"empty request": {
			r: &http.Request{},
		},
		"request with remote address": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
			},
			ip: net.IPv4(99, 99, 99, 99),
		},
		"request with xRealIP header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Real-IP": {"88.88.88.88"},
				}),
			},
			ip: net.IPv4(88, 88, 88, 88),
		},
		"request with xRealIP header and public XForwardedFor IP": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Real-IP":       {"77.77.77.77"},
					"X-Forwarded-For": {"88.88.88.88"},
				}),
			},
			ip: net.IPv4(88, 88, 88, 88),
		},
		"request with xRealIP header and private XForwardedFor IP": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Real-IP":       {"88.88.88.88"},
					"X-Forwarded-For": {"10.0.0.5"},
				}),
			},
			ip: net.IPv4(88, 88, 88, 88),
		},
		"request with single public IP in xForwardedFor header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"88.88.88.88"},
				}),
			},
			ip: net.IPv4(88, 88, 88, 88),
		},
		"request with two public IPs in xForwardedFor header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"88.88.88.88", "77.77.77.77"},
				}),
			},
			ip: net.IPv4(88, 88, 88, 88),
		},
		"request with private and public IPs in xForwardedFor header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"192.168.1.5", "88.88.88.88", "10.0.0.1", "77.77.77.77"},
				}),
			},
			ip: net.IPv4(88, 88, 88, 88),
		},
		"request with single private IP in xForwardedFor header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"192.168.1.5"},
				}),
			},
			ip: net.IPv4(192, 168, 1, 5),
		},
		"request with private IPs in xForwardedFor header": {
			r: &http.Request{
				RemoteAddr: "99.99.99.99",
				Header: makeHeader(map[string][]string{
					"X-Forwarded-For": {"192.168.1.5", "10.0.0.17"},
				}),
			},
			ip: net.IPv4(192, 168, 1, 5),
		},
	}
	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			parser := NewParser()
			ip := parser.ParseHTTPRequest(testCase.r)
			assert.Equal(t, testCase.ip, ip)
		})
	}
}
