package verification

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmailFormatNotValid  = errors.New("email format is not valid")
	ErrEmailHostUnreachable = errors.New("email host is not reachable")
)

// ValidateEmail verifies the format and the existence of an email address with a MX lookup.
func (v *Verifier) ValidateEmail(email string) error {
	if !v.Regex.MatchEmail(email) {
		return ErrEmailFormatNotValid
	}
	i := strings.LastIndexByte(email, '@')
	host := email[i+1:]
	_, err := v.mxLookup(host)
	if err != nil {
		return fmt.Errorf("%w: for host %s: %s",
			ErrEmailHostUnreachable, host, err)
	}
	return nil
}
