package verification

import (
	"fmt"
	"net"
	"strings"
)

// LookupMXFunc is the type of the net.LookupMX function
type LookupMXFunc func(name string) ([]*net.MX, error)

// ValidateEmail verifies the format and the existence of an email address with a MX lookup
func ValidateEmail(email string, mxLookup LookupMXFunc) error {
	if !MatchEmail(email) {
		return fmt.Errorf("email format of email address %q is invalid", email)
	}
	i := strings.LastIndexByte(email, '@')
	host := email[i+1:]
	_, err := mxLookup(host)
	if err != nil {
		return fmt.Errorf("host of email address %q cannot be reached: %w", email, err)
	}
	return nil
}
