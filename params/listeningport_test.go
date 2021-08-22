package params

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ListeningPort(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		listeningPort uint16
		warning       string
		err           error
	}{
		"key with valid value": {
			envValue:      "9000",
			listeningPort: 9000,
		},
		"key with valid warning value": {
			envValue:      "60000",
			listeningPort: 60000,
			warning:       "listening port 60000 is in the dynamic/private ports range (above 49151)",
		},
		"key without value": {
			envValue:      "",
			listeningPort: 0,
			err:           fmt.Errorf("port is not an integer: "),
		},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("9000")},
			listeningPort: 9000,
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           ErrNoValue,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const expectedUID = 1000
			e := &Env{
				kv: map[string]string{"key": tc.envValue},
				getuid: func() int {
					return expectedUID
				},
			}
			listeningPort, warning, err := e.ListeningPort("key", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.warning, warning)
			assert.Equal(t, tc.listeningPort, listeningPort)
		})
	}
}

func Test_checkListeningPort(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		port    uint16
		uid     int
		warning string
		err     error
	}{
		"valid port": {
			port: 9000,
		},
		"dynamic range": {
			port:    60000,
			warning: "listening port 60000 is in the dynamic/private ports range (above 49151)",
		},
		"privileged as root": {
			port:    100,
			warning: "listening port 100 allowed to be in the reserved system ports range as you are running as root",
		},
		"privileged as windows": {
			uid:     -1,
			port:    100,
			warning: "listening port 100 allowed to be in the reserved system ports range as you are running in Windows",
		},
		"privileged as non root": {
			uid:  1000,
			port: 100,
			err:  fmt.Errorf("%w: port 100", ErrReservedListeningPort),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := &Env{
				getuid: func() int {
					return tc.uid
				},
			}
			warning, err := e.checkListeningPort(tc.port)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.warning, warning)
		})
	}
}
