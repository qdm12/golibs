package verification

import (
	"fmt"
	"strings"
)

// ValidateEmail verifies the format and the existence of an email address with a MX lookup
func (v *verifier) ValidateEmail(email string) error {
	if !v.MatchEmail(email) {
		return fmt.Errorf("email format of email address %q is invalid", email)
	}
	i := strings.LastIndexByte(email, '@')
	host := email[i+1:]
	_, err := v.mxLookup(host)
	if err != nil {
		return fmt.Errorf("host of email address %q cannot be reached: %w", email, err)
	}
	return nil
}
