package verification

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VerifyPort(t *testing.T) {
	t.Parallel()
	v := NewVerifier()
	tests := map[string]struct {
		port string
		err  error
	}{
		"invalid alpha port":    {"aa", errors.New("port is not an integer: aa")},
		"invalid floating port": {"5000.55", errors.New("port is not an integer: 5000.55")},
		"invalid port 0":        {"0", errors.New("port cannot be lower than 1: 0")},
		"invalid port -1":       {"-1", errors.New("port cannot be lower than 1: -1")},
		"invalid port 70000":    {"70000", errors.New("port cannot be higher than 65535: 70000")},
		"valid port 8000":       {"8000", nil},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := v.VerifyPort(tc.port)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
