package network

import (
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_IPHeaders_isVoid(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		headers *IPHeaders
		void    bool
	}{
		"nil headers":       {nil, true},
		"empty headers":     {&IPHeaders{}, true},
		"non empty headers": {&IPHeaders{RemoteAddress: "a"}, false},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			void := tc.headers.isVoid()
			assert.Equal(t, tc.void, void)
		})
	}
}

func Test_IPHeaders_String(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		headers *IPHeaders
		s       string
	}{
		"nil headers": {s: "remoteAddr= | xRealIP= | xForwardedFor="},
		"empty headers": {
			headers: &IPHeaders{},
			s:       `remoteAddr="" | xRealIP="" | xForwardedFor=""`,
		},
		"non empty headers": {
			headers: &IPHeaders{RemoteAddress: "a", XRealIP: "bvc e"},
			s:       `remoteAddr="a" | xRealIP="bvc e" | xForwardedFor=""`,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := tc.headers.String()
			assert.Equal(t, tc.s, s)
		})
	}
}

func Test_GetClientIPHeaders(t *testing.T) {
	t.Parallel()
	emptyLogging, _ := logging.NewEmptyLogger()
	m := NewIPManager(emptyLogging)
	tests := map[string]struct {
		request *http.Request
		headers IPHeaders
	}{
		"nil request":   {nil, IPHeaders{}},
		"empty request": {&http.Request{}, IPHeaders{}},
		"normal request": {&http.Request{
			RemoteAddr: " a ",
			Header: http.Header{
				"X-Real-Ip":       []string{" b"},
				"X-Forwarded-For": []string{"c"},
			},
		}, IPHeaders{
			RemoteAddress: "a",
			XRealIP:       "b",
			XForwardedFor: "c",
		}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			headers := m.GetClientIPHeaders(tc.request)
			assert.Equal(t, tc.headers, headers)
		})
	}
}

func Test_GetClientIP(t *testing.T) {
	t.Parallel()
	emptyLogging, _ := logging.NewEmptyLogger()
	m := NewIPManager(emptyLogging)
	tests := map[string]struct {
		request *http.Request
		IP      string
		err     error
	}{
		"nil request": {
			err: fmt.Errorf("no IP address found in client request"),
		},
		"empty request": {
			request: &http.Request{},
			err:     fmt.Errorf("no IP address found in client request"),
		},
		"request with remote address": {
			request: &http.Request{RemoteAddr: "54.54.54.54:8888"},
			IP:      "54.54.54.54",
		},
		"request with bad remote address": {
			request: &http.Request{RemoteAddr: "54.54.54.300:8888"},
			err:     fmt.Errorf(`IP address "54.54.54.300" is not valid`),
		},
		"request with remote address and 1 xRealIP": {
			request: &http.Request{
				RemoteAddr: "54.54.54.54:8888",
				Header:     http.Header{"X-Real-Ip": []string{"55.55.55.55"}, "X-Forwarded-For": []string{""}},
			},
			IP: "55.55.55.55",
		},
		"request with remote address, 1 xRealIP and 2 XForwardedFor": {
			request: &http.Request{
				RemoteAddr: "54.54.54.54:8888",
				Header: http.Header{
					"X-Real-Ip":       []string{"55.55.55.55"},
					"X-Forwarded-For": []string{"192.168.1.10 , 10.0.0.1"},
				},
			},
			IP: "55.55.55.55",
		},
		"request with remote address and 2 xForwardedFor": {
			request: &http.Request{
				RemoteAddr: "54.54.54.54:8888",
				Header:     http.Header{"X-Real-Ip": []string{""}, "X-Forwarded-For": []string{"192.168.1.10 , 10.0.0.1"}},
			},
			IP: "10.0.0.1",
		},
		"request with remote address, 1 XRealIP and 3 XForwardedFor": {
			request: &http.Request{
				RemoteAddr: "54.54.54.54:8888",
				Header: http.Header{
					"X-Real-Ip":       []string{"55.55.55.55"},
					"X-Forwarded-For": []string{"192.168.1.10 , 10.0.0.1, 56.56.56.56"},
				},
			},
			IP: "56.56.56.56",
		},
		"request with remote address, and bad XForwardedFor": {
			request: &http.Request{
				RemoteAddr: "54.54.54.54:8888",
				Header:     http.Header{"X-Forwarded-For": []string{"x"}},
			},
			IP: "54.54.54.54",
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			IP, err := m.GetClientIP(tc.request)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.IP, IP)
		})
	}
}

func Test_ipIsValid(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		IP    string
		valid bool
	}{
		"no address":           {"", false},
		"bad digit address":    {"125", false},
		"bad text address":     {"bla", false},
		"simple IPv4 address":  {"1.1.0.1", true},
		"simple IPv6 address":  {"::0", true},
		"complex IPv4 address": {"192.168.25.218", true},
		"complex IPv6 address": {"fdf7:8fb3:2a0:62d:0:0:0:0", true},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			valid := ipIsValid(tc.IP)
			assert.Equal(t, tc.valid, valid)
		})
	}
}

func Test_netIPIsPrivate(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		netIP   net.IP
		private bool
	}{
		"no IP":        {nil, false},
		"Public IPv4":  {net.ParseIP("15.1.25.1"), false},
		"Public IPv6":  {net.ParseIP("2001:db8::8a2e:370:7334"), false},
		"Private IPv4": {net.ParseIP("192.168.25.218"), true},
		"Private IPv6": {net.ParseIP("fd8d:8d72:b629:0f87:0000:0000:0000:0000"), true},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			private := netIPIsPrivate(tc.netIP)
			assert.Equal(t, tc.private, private)
		})
	}
}

func Test_splitHostPort(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		address string
		IP      string
		port    string
		err     error
	}{
		"address with port only": {
			address: ":8000",
			port:    "8000",
		},
		"address with IPv4 address and empty port": {
			address: "192.165.100.100:",
			IP:      "192.165.100.100",
		},
		"address with IPv4 address and port": {
			address: "192.165.100.100:8000",
			IP:      "192.165.100.100",
			port:    "8000",
		},
		"address with IPv4 address and too many colons": {
			address: "192.165.100.100:8000:500",
			err:     fmt.Errorf("address 192.165.100.100:8000:500: too many colons in address"),
		},
		"address with IPv6 address and empty port": {
			address: "2001:db8::8a2e:370:7334:",
			IP:      "2001:db8::8a2e:370:7334",
		},
		"address with IPv6 address and port": {
			address: "2001:db8::8a2e:370:7334:8000",
			IP:      "2001:db8::8a2e:370:7334",
			port:    "8000",
		},
		"address with IPv6 address and too many colons": {
			address: "2001:db8::8a2e:370:511:32:250:7334:8000",
			err:     fmt.Errorf("address 2001:db8::8a2e:370:511:32:250:7334:8000: too many colons in address"),
		},
		"address with [IPv6 address] and empty port": {
			address: "[2001:db8::8a2e:370:7334]:",
			IP:      "2001:db8::8a2e:370:7334",
		},
		"address with [IPv6 address] and port": {
			address: "[2001:db8::8a2e:370:7334]:8000",
			IP:      "2001:db8::8a2e:370:7334",
			port:    "8000",
		},
		"address with [IPv6 address] and too many colons after port": {
			address: "[2001:db8::8a2e:370:7334]:8000:500",
			err:     fmt.Errorf("address [2001:db8::8a2e:370:7334]:8000:500: too many colons in address"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			IP, port, err := splitHostPort(tc.address)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.IP, IP)
			assert.Equal(t, tc.port, port)
		})
	}
}

func Test_getRemoteIP(t *testing.T) {
	t.Parallel()
	verifyPort := verification.NewVerifier().VerifyPort
	tests := map[string]struct {
		remoteAddr string
		IP         string
		err        error
	}{
		"address with invalid IPv4 address and no port": {
			remoteAddr: "192.165.300.100",
			err:        fmt.Errorf(`IP address "192.165.300.100" is not valid`),
		},
		"address with invalid IPv4 address and valid port": {
			remoteAddr: "192.165.300.100:8000",
			err:        fmt.Errorf(`IP address "192.165.300.100" is not valid`),
		},
		"address with valid IPv4 address and invalid port": {
			remoteAddr: "192.165.200.100:70000",
			err:        fmt.Errorf(`remote address "192.165.200.100:70000" is invalid: port 70000 cannot be higher than 65535`),
		},
		"address with too many colons": {
			remoteAddr: "192.165.200.100:8000:500",
			err:        fmt.Errorf("address 192.165.200.100:8000:500: too many colons in address"),
		},
		"address with valid IPv4 address and no port": {
			remoteAddr: "192.165.168.100",
			IP:         "192.165.168.100",
		},
		"address with valid IPv4 address and empty port": {
			remoteAddr: "192.165.168.100:",
			IP:         "192.165.168.100",
		},
		"address with valid IPv4 address and valid port": {
			remoteAddr: "192.165.168.100:8000",
			IP:         "192.165.168.100",
		},
		"address with valid IPv4 address and invalid port with letters": {
			remoteAddr: "192.165.168.100:80a0b",
			err:        fmt.Errorf(`remote address "192.165.168.100:80a0b" is invalid: port "80a0b" is not a valid integer`),
		},
		"address with valid IPv6 address and no port": {
			remoteAddr: "2001:db8::8a2e:370:7334",
			IP:         "2001:db8::8a2e:370",
		},
		"address with valid IPv6 address and empty port": {
			remoteAddr: "2001:db8::8a2e:370:7334:",
			IP:         "2001:db8::8a2e:370:7334",
		},
		"address with valid IPv6 address and valid port": {
			remoteAddr: "2001:db8::8a2e:370:7334:8000",
			IP:         "2001:db8::8a2e:370:7334",
		},
		"address with [valid IPv6 address] and valid port": {
			remoteAddr: "[2001:db8::8a2e:370:7334]:8000",
			IP:         "2001:db8::8a2e:370:7334",
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			IP, err := getRemoteIP(verifyPort, tc.remoteAddr)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.IP, IP)
		})
	}
}

func Test_extractPublicIPs(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		IPs       []string
		publicIPs []string
		warnings  []string
	}{
		"no IPs": {IPs: []string{}},
		"2 private IPv4": {
			IPs: []string{"127.0.0.3", "192.168.178.5"},
		},
		"2 private and 1 invalid IPv4": {
			IPs:      []string{"127.0.0.3", "192.168.178.5", "58"},
			warnings: []string{`IP address "58" is not valid`},
		},
		"2 private and 1 public IPv4": {
			IPs:       []string{"127.0.0.3", "192.168.178.5", "58.58.58.58"},
			publicIPs: []string{"58.58.58.58"},
		},
		"1 private and 1 public IPv6": {
			IPs:       []string{"fd8d:8d72:b629:0f87:0000:0000:0000:0000", "2001:db8::8a2e:370:7334"},
			publicIPs: []string{"2001:db8::8a2e:370:7334"},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			publicIPs, warnings := extractPublicIPs(tc.IPs)
			assert.ElementsMatch(t, tc.publicIPs, publicIPs)
			assert.ElementsMatch(t, tc.warnings, warnings)
		})
	}
}

func Test_getXForwardedIPs(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		XForwardedFor string
		IPs           []string
		warnings      []string
	}{
		"no XForwardedFor": {},
		"3 IPv4": {
			XForwardedFor: "192.99.99.58, 54.54.55.56,54.54.55.100",
			IPs:           []string{"192.99.99.58", "54.54.55.56", "54.54.55.100"},
		},
		"2 IPv4 and 1 invalid with multiple spaces": {
			XForwardedFor: "192.99.99.58,  54.54.55.56,99.99.87",
			IPs:           []string{"192.99.99.58", "54.54.55.56"},
			warnings:      []string{`IP address "99.99.87" is not valid`},
		},
		"1 IPv4 and 1 IPv6 with \\t": {
			XForwardedFor: " 192.99.99.58, \t2001:db8::8a2e:370:7334",
			IPs:           []string{"192.99.99.58", "2001:db8::8a2e:370:7334"},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			IPs, warnings := getXForwardedIPs(tc.XForwardedFor)
			assert.ElementsMatch(t, tc.IPs, IPs)
			assert.ElementsMatch(t, tc.warnings, warnings)
		})
	}
}
