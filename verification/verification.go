package verification

import "net"

type Verifier interface {
	// VerifyPort verifies a port number is valid
	VerifyPort(port string) error
	// ValidateEmail verifies an email format is valid and performs a DNS lookup
	// to verify the email address does exist
	ValidateEmail(email string) error
	// SearchIPv4 extracts all IPv4 addresses found in the given string
	SearchIPv4(s string) []string
	// SearchIPv6 extracts all IPv6 addresses found in the given string
	SearchIPv6(s string) []string
	// SearchEmail extracts all email addresses found in the given string
	SearchEmail(s string) []string
	// SearchPhone extracts all phone numbers found in the given string
	SearchPhone(s string) []string
	MatchEmail(s string) bool
	MatchPhoneIntl(s string) bool
	MatchPhoneLocal(s string) bool
	MatchDomain(s string) bool
	MatchHostname(s string) bool
	MatchRootURL(s string) bool
	MatchMD5String(s string) bool
	Match64BytesHex(s string) bool
}

type verifier struct {
	mxLookup func(name string) ([]*net.MX, error)
	Regex
}

func NewVerifier() Verifier {
	return &verifier{
		mxLookup: net.LookupMX,
		Regex:    NewRegex(),
	}
}
