package params

import (
	"errors"
	"testing"

	"github.com/qdm12/golibs/verification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_verifyListeningPort(t *testing.T) {
	t.Parallel()
	v := verification.NewVerifier()
	tests := map[string]struct {
		listeningPort string
		uid           int
		warning       string
		err           error
	}{
		"invalid port":                    {"a", 0, "", errors.New("listening port: port \"a\" is not a valid integer")},
		"reserved system port as root":    {"100", 0, "listening port 100 allowed to be in the reserved system ports range as you are running as root", nil},
		"reserved system port as Windows": {"100", -1, "listening port 100 allowed to be in the reserved system ports range as you are running in Windows", nil},
		"reserved system port as UID > 0": {"100", 1000, "", errors.New("listening port 100 cannot be in the reserved system ports range (1 to 1023) when running without root")},
		"dynamic/private port":            {"50000", 0, "listening port 50000 is in the dynamic/private ports range (above 49151)", nil},
		"valid port":                      {"8000", 1000, "", nil},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			warning, err := verifyListeningPort(v, tc.listeningPort, tc.uid)
			assert.Equal(t, tc.warning, warning)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
