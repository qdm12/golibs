package params

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Path(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		fpAbs         func(p string) (string, error)
		path          string
		err           error
	}{
		"valid path": {
			envValue: "/path",
			fpAbs: func(p string) (string, error) {
				return "/real/path", nil
			},
			path: "/real/path",
		},
		"using filepath.Abs": {
			envValue: "/path",
			fpAbs:    filepath.Abs,
			path:     "/path",
		},
		"get error": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           ErrNoValue,
		},
		"abs error": {
			envValue: "/path",
			fpAbs: func(p string) (string, error) {
				return "", errors.New("abs error")
			},
			err: errors.New(`invalid filepath: : abs error`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &Env{
				kv:    map[string]string{"key": tc.envValue},
				fpAbs: tc.fpAbs,
			}
			path, err := e.Path("key", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.path, path)
		})
	}
}
