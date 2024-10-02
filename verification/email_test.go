package verification

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ValidateEmail(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		email      string
		mxLookup   func(name string) ([]*net.MX, error)
		errWrapped error
		errMessage string
	}{
		"Invalid email format": {
			email:      "aa",
			errWrapped: ErrEmailFormatNotValid,
			errMessage: "email format is not valid: aa",
		},
		"Valid email format but not existing": {
			email: "aa@aa.aa",
			mxLookup: func(_ string) ([]*net.MX, error) {
				return nil, errors.New("not existing")
			},
			errWrapped: ErrEmailHostUnreachable,
			errMessage: "email host is not reachable: for host aa.aa: not existing",
		},
		"Valid email format and existing": {
			email: "aa@aa.aa",
			mxLookup: func(_ string) ([]*net.MX, error) {
				return nil, nil
			},
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			v := &Verifier{mxLookup: tc.mxLookup, Regex: NewRegex()}
			err := v.ValidateEmail(tc.email)
			assert.ErrorIs(t, err, tc.errWrapped)
			if tc.errWrapped != nil {
				assert.EqualError(t, err, tc.errMessage)
			}
		})
	}
}
