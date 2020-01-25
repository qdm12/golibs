package verification

import "net"

type Verifier interface {
	VerifyPort(port string) error
	ValidateEmail(email string) error
	SearchIPv4(s string) []string
	SearchIPv6(s string) []string
	SearchEmail(s string) []string
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
}

func NewVerifier() Verifier {
	return &verifier{
		mxLookup: net.LookupMX,
	}
}
