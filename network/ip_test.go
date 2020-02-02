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
		"nil headers":       {nil, "remoteAddr= | xRealIP= | xForwardedFor="},
		"empty headers":     {&IPHeaders{}, "remoteAddr=\"\" | xRealIP=\"\" | xForwardedFor=\"\""},
		"non empty headers": {&IPHeaders{RemoteAddress: "a", XRealIP: "bvc e"}, "remoteAddr=\"a\" | xRealIP=\"bvc e\" | xForwardedFor=\"\""},
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
		"nil request":                     {nil, "", fmt.Errorf("no IP address found in client request")},
		"empty request":                   {&http.Request{}, "", fmt.Errorf("no IP address found in client request")},
		"request with remote address":     {&http.Request{RemoteAddr: "54.54.54.54:8888"}, "54.54.54.54", nil},
		"request with bad remote address": {&http.Request{RemoteAddr: "54.54.54.300:8888"}, "", fmt.Errorf("IP address \"54.54.54.300\" is not valid")},
		"request with remote address and 1 xRealIP": {
			&http.Request{
				RemoteAddr: "54.54.54.54:8888",
				Header:     http.Header{"X-Real-Ip": []string{"55.55.55.55"}, "X-Forwarded-For": []string{""}},
			},
			"55.55.55.55",
			nil,
		},
		"request with remote address, 1 xRealIP and 2 XForwardedFor": {
			&http.Request{
				RemoteAddr: "54.54.54.54:8888",
				Header:     http.Header{"X-Real-Ip": []string{"55.55.55.55"}, "X-Forwarded-For": []string{"192.168.1.10 , 10.0.0.1"}},
			},
			"55.55.55.55",
			nil,
		},
		"request with remote address and 2 xForwardedFor": {
			&http.Request{
				RemoteAddr: "54.54.54.54:8888",
				Header:     http.Header{"X-Real-Ip": []string{""}, "X-Forwarded-For": []string{"192.168.1.10 , 10.0.0.1"}},
			},
			"10.0.0.1",
			nil,
		},
		"request with remote address, 1 XRealIP and 3 XForwardedFor": {
			&http.Request{
				RemoteAddr: "54.54.54.54:8888",
				Header:     http.Header{"X-Real-Ip": []string{"55.55.55.55"}, "X-Forwarded-For": []string{"192.168.1.10 , 10.0.0.1, 56.56.56.56"}},
			},
			"56.56.56.56",
			nil,
		},
		"request with remote address, and bad XForwardedFor": {
			&http.Request{
				RemoteAddr: "54.54.54.54:8888",
				Header:     http.Header{"X-Forwarded-For": []string{"x"}},
			},
			"54.54.54.54",
			nil,
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
		"address with port only":                                     {":8000", "", "8000", nil},
		"address with IPv4 address and empty port":                   {"192.165.100.100:", "192.165.100.100", "", nil},
		"address with IPv4 address and port":                         {"192.165.100.100:8000", "192.165.100.100", "8000", nil},
		"address with IPv4 address and too many colons":              {"192.165.100.100:8000:500", "", "", fmt.Errorf("address 192.165.100.100:8000:500: too many colons in address")},
		"address with IPv6 address and empty port":                   {"2001:db8::8a2e:370:7334:", "2001:db8::8a2e:370:7334", "", nil},
		"address with IPv6 address and port":                         {"2001:db8::8a2e:370:7334:8000", "2001:db8::8a2e:370:7334", "8000", nil},
		"address with IPv6 address and too many colons":              {"2001:db8::8a2e:370:511:32:250:7334:8000", "", "", fmt.Errorf("address 2001:db8::8a2e:370:511:32:250:7334:8000: too many colons in address")},
		"address with [IPv6 address] and empty port":                 {"[2001:db8::8a2e:370:7334]:", "2001:db8::8a2e:370:7334", "", nil},
		"address with [IPv6 address] and port":                       {"[2001:db8::8a2e:370:7334]:8000", "2001:db8::8a2e:370:7334", "8000", nil},
		"address with [IPv6 address] and too many colons after port": {"[2001:db8::8a2e:370:7334]:8000:500", "", "", fmt.Errorf("address [2001:db8::8a2e:370:7334]:8000:500: too many colons in address")},
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
		"address with invalid IPv4 address and no port":                 {"192.165.300.100", "", fmt.Errorf("IP address \"192.165.300.100\" is not valid")},
		"address with invalid IPv4 address and valid port":              {"192.165.300.100:8000", "", fmt.Errorf("IP address \"192.165.300.100\" is not valid")},
		"address with valid IPv4 address and invalid port":              {"192.165.200.100:70000", "", fmt.Errorf("remote address \"192.165.200.100:70000\" is invalid: port 70000 cannot be higher than 65535")},
		"address with too many colons":                                  {"192.165.200.100:8000:500", "", fmt.Errorf("address 192.165.200.100:8000:500: too many colons in address")},
		"address with valid IPv4 address and no port":                   {"192.165.168.100", "192.165.168.100", nil},
		"address with valid IPv4 address and empty port":                {"192.165.168.100:", "192.165.168.100", nil},
		"address with valid IPv4 address and valid port":                {"192.165.168.100:8000", "192.165.168.100", nil},
		"address with valid IPv4 address and invalid port with letters": {"192.165.168.100:80a0b", "", fmt.Errorf("remote address \"192.165.168.100:80a0b\" is invalid: port \"80a0b\" is not a valid integer")},
		"address with valid IPv6 address and no port":                   {"2001:db8::8a2e:370:7334", "2001:db8::8a2e:370", nil},
		"address with valid IPv6 address and empty port":                {"2001:db8::8a2e:370:7334:", "2001:db8::8a2e:370:7334", nil},
		"address with valid IPv6 address and valid port":                {"2001:db8::8a2e:370:7334:8000", "2001:db8::8a2e:370:7334", nil},
		"address with [valid IPv6 address] and valid port":              {"[2001:db8::8a2e:370:7334]:8000", "2001:db8::8a2e:370:7334", nil},
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
		"no IPs":                       {[]string{}, nil, nil},
		"2 private IPv4":               {[]string{"127.0.0.3", "192.168.178.5"}, nil, nil},
		"2 private and 1 invalid IPv4": {[]string{"127.0.0.3", "192.168.178.5", "58"}, nil, []string{`IP address "58" is not valid`}},
		"2 private and 1 public IPv4":  {[]string{"127.0.0.3", "192.168.178.5", "58.58.58.58"}, []string{"58.58.58.58"}, nil},
		"1 private and 1 public IPv6":  {[]string{"fd8d:8d72:b629:0f87:0000:0000:0000:0000", "2001:db8::8a2e:370:7334"}, []string{"2001:db8::8a2e:370:7334"}, nil},
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
		"no XForwardedFor": {"", nil, nil},
		"3 IPv4":           {"192.99.99.58, 54.54.55.56,54.54.55.100", []string{"192.99.99.58", "54.54.55.56", "54.54.55.100"}, nil},
		"2 IPv4 and 1 invalid with multiple spaces": {"192.99.99.58,  54.54.55.56,99.99.87", []string{"192.99.99.58", "54.54.55.56"}, []string{`IP address "99.99.87" is not valid`}},
		"1 IPv4 and 1 IPv6 with \\t":                {" 192.99.99.58, \t2001:db8::8a2e:370:7334", []string{"192.99.99.58", "2001:db8::8a2e:370:7334"}, nil},
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
