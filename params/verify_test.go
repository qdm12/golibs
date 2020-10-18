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
		"invalid port": {
			listeningPort: "a",
			uid:           0,
			err:           errors.New(`listening port: port "a" is not a valid integer`),
		},
		"reserved system port as root": {
			listeningPort: "100",
			uid:           0,
			warning:       "listening port 100 allowed to be in the reserved system ports range as you are running as root",
		},
		"reserved system port as Windows": {
			listeningPort: "100",
			uid:           -1,
			warning:       "listening port 100 allowed to be in the reserved system ports range as you are running in Windows",
		},
		"reserved system port as UID > 0": {
			listeningPort: "100",
			uid:           1000,
			err:           errors.New("listening port 100 cannot be in the reserved system ports range (1 to 1023) when running without root"), //nolint:lll
		},
		"dynamic/private port": {
			listeningPort: "50000",
			uid:           0,
			warning:       "listening port 50000 is in the dynamic/private ports range (above 49151)",
		},
		"valid port": {
			listeningPort: "8000",
			uid:           1000,
		},
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
