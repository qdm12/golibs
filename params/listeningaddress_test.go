package params

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ListeningAddress(t *testing.T) {
	t.Parallel()
	const key = "LISTENING_ADDRESS"
	tests := map[string]struct {
		envValue string
		options  []OptionSetter
		address  string
		warning  string
		err      error
	}{
		"success": {
			envValue: "0.0.0.0:8000",
			address:  "0.0.0.0:8000",
		},
		"env get error": {
			envValue: "",
			options:  []OptionSetter{Compulsory()},
			err:      ErrNoValue,
		},
		"split host port error": {
			envValue: "0.0.0.0",
			err:      errors.New("address 0.0.0.0: missing port in address"),
		},
		"bad port string": {
			envValue: "0.0.0.0:a",
			err:      errors.New(`invalid port: strconv.Atoi: parsing "a": invalid syntax`),
		},
		"reserved port error": {
			envValue: "0.0.0.0:100",
			err:      errors.New(`invalid port: listening port cannot be in the reserved system ports range (1 to 1023) when running without root: port 100`), //nolint:lll
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			e := &Env{
				getenv: func(key string) string {
					return tc.envValue
				},
				getuid: func() int {
					const uid = 1000
					return uid
				},
			}

			address, warning, err := e.ListeningAddress(key, tc.options...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.warning, warning)
			assert.Equal(t, tc.address, address)
		})
	}
}
