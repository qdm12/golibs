package verification

import (
	"errors"
	"testing"

	"github.com/qdm12/golibs/helpers"
)

func Test_VerifyPort(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		port string
		err  error
	}{
		"invalid alpha port":    {"aa", errors.New("port \"aa\" is not a valid integer")},
		"invalid floating port": {"5000.55", errors.New("port \"5000.55\" is not a valid integer")},
		"invalid port 0":        {"0", errors.New("port 0 cannot be lower than 1")},
		"invalid port -1":       {"-1", errors.New("port -1 cannot be lower than 1")},
		"invalid port 70000":    {"70000", errors.New("port 70000 cannot be higher than 65535")},
		"valid port 8000":       {"8000", nil},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := VerifyPort(tc.port)
			helpers.AssertErrorsEqual(t, tc.err, err)
		})
	}
}
