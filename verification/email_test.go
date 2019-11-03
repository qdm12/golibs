package verification

import (
	"errors"
	"net"
	"testing"

	"github.com/qdm12/golibs/helpers"
)

func Test_ValidateEmail(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		email    string
		mxLookup LookupMXFunc
		err      error
	}{
		"Invalid email format": {
			"aa",
			nil,
			errors.New("email format of email address \"aa\" is invalid"),
		},
		"Valid email format but not existing": {
			"aa@aa.aa",
			func(name string) ([]*net.MX, error) {
				return nil, errors.New("not existing")
			},
			errors.New("host of email address \"aa@aa.aa\" cannot be reached: not existing"),
		},
		"Valid email format and existing": {
			"aa@aa.aa",
			func(name string) ([]*net.MX, error) {
				return nil, nil
			},
			nil,
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			err := ValidateEmail(tc.email, tc.mxLookup)
			helpers.AssertErrorsEqual(t, tc.err, err)
		})
	}
}
