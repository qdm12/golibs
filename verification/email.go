package verification

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmailFormatNotValid  = errors.New("email format is not valid")
	ErrEmailHostUnreachable = errors.New("email host is not reachable")
)

// ValidateEmail verifies the format and the existence of an email address with a MX lookup.
func ValidateEmail(ctx context.Context, email string,
	mxLookuper MXLookuper) error {
	if !MatchEmail(email) {
		return fmt.Errorf("%w: %s", ErrEmailFormatNotValid, email)
	}

	i := strings.LastIndexByte(email, '@')
	host := email[i+1:]
	records, err := mxLookuper.LookupMX(ctx, host)

	switch {
	case err != nil:
		return fmt.Errorf("%w: for host %s: %w",
			ErrEmailHostUnreachable, host, err)
	case len(records) == 0:
		return fmt.Errorf("%w: for host %s: no MX record found",
			ErrEmailHostUnreachable, host)
	default:
		return nil
	}
}
