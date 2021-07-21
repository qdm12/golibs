package verification

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ValidateEmail(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		email    string
		mxLookup func(name string) ([]*net.MX, error)
		err      error
	}{
		"Invalid email format": {
			"aa",
			nil,
			ErrEmailFormatNotValid,
		},
		"Valid email format but not existing": {
			"aa@aa.aa",
			func(name string) ([]*net.MX, error) {
				return nil, errors.New("not existing")
			},
			errors.New("email host is not reachable: for host aa.aa: not existing"),
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
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			v := &verifier{mxLookup: tc.mxLookup, Regex: NewRegex()}
			err := v.ValidateEmail(tc.email)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
