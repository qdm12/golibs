package network

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewConnectivity(t *testing.T) {
	t.Parallel()
	c := NewConnectivity(time.Second)
	assert.NotNil(t, c)
}

func Test_ConnectivityChecks(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		domains []string
		DNSErr  error
		errs    []error
	}{
		"no error":     {nil, nil, nil},
		"error for 1":  {[]string{"domain.com"}, fmt.Errorf("error"), []error{fmt.Errorf("Domain name resolution is not working for domain.com: error")}},
		"errors for 2": {[]string{"domain.com", "domain2.com"}, fmt.Errorf("error"), []error{fmt.Errorf("Domain name resolution is not working for domain.com: error"), fmt.Errorf("Domain name resolution is not working for domain2.com: error")}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			checkDNS := func(domain string) error { return tc.DNSErr }
			mockClient := &MockClient{}
			for _, domain := range tc.domains {
				mockClient.On("GetContent", "http://"+domain).
					Return(nil, 200, nil).Once()
				mockClient.On("GetContent", "https://"+domain).
					Return(nil, 200, nil).Once()
			}
			connectivity := &connectivity{
				checkDNS: checkDNS,
				client:   mockClient,
			}
			errs := connectivity.Checks(tc.domains...)
			expectedErrsString := []string{}
			for _, err := range tc.errs {
				expectedErrsString = append(expectedErrsString, err.Error())
			}
			errsString := []string{}
			for _, err := range errs {
				errsString = append(errsString, err.Error())
			}
			assert.ElementsMatch(t, expectedErrsString, errsString)
			mockClient.AssertExpectations(t)
		})
	}
}

func Test_connectivityCheck(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		DNSErr   error
		HTTPErr  error
		HTTPSErr error
		errs     []error
	}{
		"no error":     {nil, nil, nil, nil},
		"DNS error":    {fmt.Errorf("error"), nil, nil, []error{fmt.Errorf("Domain name resolution is not working for domain.com: error")}},
		"HTTP error":   {nil, fmt.Errorf("error"), nil, []error{fmt.Errorf("HTTP GET failed for http://domain.com: error")}},
		"HTTPS error":  {nil, nil, fmt.Errorf("error"), []error{fmt.Errorf("HTTP GET failed for https://domain.com: error")}},
		"Mixed errors": {fmt.Errorf("error"), nil, fmt.Errorf("error"), []error{fmt.Errorf("Domain name resolution is not working for domain.com: error"), fmt.Errorf("HTTP GET failed for https://domain.com: error")}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			checkDNS := func(domain string) error { return tc.DNSErr }
			mockClient := &MockClient{}
			mockClient.On("GetContent", "http://domain.com").
				Return(nil, 200, tc.HTTPErr).Once()
			mockClient.On("GetContent", "https://domain.com").
				Return(nil, 200, tc.HTTPSErr).Once()
			errs := connectivityCheck("domain.com", checkDNS, mockClient)
			expectedErrsString := []string{}
			for _, err := range tc.errs {
				expectedErrsString = append(expectedErrsString, err.Error())
			}
			errsString := []string{}
			for _, err := range errs {
				errsString = append(errsString, err.Error())
			}
			assert.ElementsMatch(t, expectedErrsString, errsString)
			mockClient.AssertExpectations(t)
		})
	}
}

func Test_httpGetCheck(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		getContentStatus int
		getContentErr    error
		err              error
	}{
		"no error":   {200, nil, nil},
		"bad status": {400, nil, fmt.Errorf("HTTP GET failed for https://domain.com: HTTP Status 400")},
		"error":      {0, fmt.Errorf("error"), fmt.Errorf("HTTP GET failed for https://domain.com: error")},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			mockClient := &MockClient{}
			mockClient.On("GetContent", "https://domain.com").
				Return(nil, tc.getContentStatus, tc.getContentErr).Once()
			err := httpGetCheck("https://domain.com", mockClient)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			mockClient.AssertExpectations(t)
		})
	}
}

func Test_domainNameResolutionCheck(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		checkDNSErr error
		err         error
	}{
		"no error": {nil, nil},
		"error":    {fmt.Errorf("error"), fmt.Errorf("Domain name resolution is not working for domain.com: error")},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			checkDNS := func(domain string) error { return tc.checkDNSErr }
			err := domainNameResolutionCheck("domain.com", checkDNS)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
