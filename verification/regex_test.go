package verification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SearchIPv4(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s     string
		finds []string
	}{
		"Find nothing in empty string":  {"", nil},
		"Find nothing":                  {"dsadsa 232.323 s", nil},
		"Find exactly":                  {"192.168.1.5", []string{"192.168.1.5"}},
		"Find multiple in text":         {"sd 192.168.1.5 1.5 1.3.5.4", []string{"192.168.1.5", "1.3.5.4"}},
		"Find in other text":            {"bla 192.168.1.0 bla", []string{"192.168.1.0"}},
		"Find in longer than normal IP": {"0.0.0.0.0", []string{"0.0.0.0"}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := SearchIPv4(tc.s)
			assert.ElementsMatch(t, tc.finds, out)
		})
	}
}

func Test_SearchIPv6(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s     string
		finds []string
	}{
		"Find nothing in empty string": {"", nil},
		"Find nothing":                 {"dsadsa 232.323 s", nil},
		"Find nothing in IPv4 address": {"192.168.1.5", nil},
		"Find exactly ::1":             {"::1", []string{"::1"}},
		"Find multiple in text":        {"2001:0db8:85a3:0000:0000:8a2e:0370:7334 sdas ::1", []string{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", "::1"}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := SearchIPv6(tc.s)
			assert.ElementsMatch(t, tc.finds, out)
		})
	}
}

func Test_SearchEmail(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s     string
		finds []string
	}{
		"Find nothing in empty string": {"", nil},
		"Find single email in text":    {"bla bla bla@bla bla@bla.co.uk bla.com", []string{"bla@bla.co.uk"}},
		"Find two emails in text":      {"bla@aa.aa bla bla@bla bla@bla.co.uk bla.com", []string{"bla@aa.aa", "bla@bla.co.uk"}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := SearchEmail(tc.s)
			assert.ElementsMatch(t, tc.finds, out)
		})
	}
}

func Test_SearchPhone(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s     string
		finds []string
	}{
		"Find nothing in empty string":                                      {"", nil},
		"Find nothing in text":                                              {"aa", nil},
		"Find international number without + sign":                          {"35226440600", []string{"35226440600"}},
		"Find international number with + sign":                             {"+35226440600", []string{"+35226440600"}},
		"Find international number with + sign and 1 space":                 {"+352 26440600", []string{"+352 26440600"}},
		"Find international number with + sign and multiple spaces in text": {"bla +1 3474 50256 2 blaaa 234", []string{"+1 3474 50256 2"}},
		"Complex case": {"b +1 3474 50256 2 fdfd 332 23d 45787e 35226440600", []string{"+1 3474 50256 2", "35226440600"}},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := SearchPhone(tc.s)
			assert.ElementsMatch(t, tc.finds, out)
		})
	}
}

func Test_MatchEmail(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s     string
		match bool
	}{
		"No match in empty string":                {"", false},
		"No match in text":                        {"aa", false},
		"No match in email without hostname":      {"aa@.aa", false},
		"No match in email without tld":           {"aa@aa", false},
		"No match in email without user":          {"@aa.aa", false},
		"match for simple email":                  {"aa@aa.aa", true},
		"no match for email made of numbers only": {"125@125.12", false},
		"match for email with numbers except tld": {"125@125.aa", true},
		"match for complex email":                 {"aaabc@aaabc.co.uk", true},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := MatchEmail(tc.s)
			assert.Equal(t, tc.match, out)
		})
	}
}

func Test_MatchPhoneIntl(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s     string
		match bool
	}{
		"no match in empty string":            {"", false},
		"no match in text":                    {"aa", false},
		"no match with number without +":      {"35226440600", false},
		"match 1":                             {"+35226440600", true},
		"no match because of space in number": {"+352 26440600", false},
		"match 2":                             {"+13474502562", true},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := MatchPhoneIntl(tc.s)
			assert.Equal(t, tc.match, out)
		})
	}
}

func Test_MatchPhoneLocal(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s     string
		match bool
	}{
		"No match in empty string":                          {"", false},
		"No match in text":                                  {"aa", false},
		"Match long number":                                 {"35226440600", true},
		"Match short number":                                {"26440600", true},
		"No match for too short number":                     {"2222", false},
		"No match for too long number":                      {"222222222222222", false},
		"No match for international number with leading +":  {"+35226440600", false},
		"No match for international number with leading 00": {"0035226440600", false},
		"match for long number with one leading 0":          {"0535226440600", true},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := MatchPhoneLocal(tc.s)
			assert.Equal(t, tc.match, out)
		})
	}
}

func Test_MatchDomain(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s     string
		match bool
	}{
		"No match in empty string":      {"", false},
		"No match in text":              {"aa", false},
		"No match as only TLD":          {".com", false},
		"No match as only composed TLD": {".co.uk", false},
		"Match":                         {"aa.aa", true},
		"Match one letter":              {"d.com", true},
		"Math numbers":                  {"765.fr", true},
		"Match composed TLD":            {"hey.co.uk", true},
		"Match composed EDU TLD":        {"nyu.ac.edu", true},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := MatchDomain(tc.s)
			assert.Equal(t, tc.match, out)
		})
	}
}

func Test_MatchHostname(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		s     string
		match bool
	}{
		"No match in empty string":      {"", false},
		"Simple one word hostname":      {"aa", true},
		"No match as only TLD":          {".com", false},
		"No match as only composed TLD": {".co.uk", false},
		"Match":                         {"aa.aa", true},
		"Match one letter":              {"d.com", true},
		"Math numbers":                  {"765.fr", true},
		"Match composed TLD":            {"hey.co.uk", true},
		"Match composed EDU TLD":        {"nyu.ac.edu", true},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			out := MatchHostname(tc.s)
			assert.Equal(t, tc.match, out)
		})
	}
}
