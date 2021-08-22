package params

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_GetInt(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		n             int
		err           error
	}{
		"key with int value": {envValue: "0"},
		"key with float value": {
			envValue: "0.00",
			err:      fmt.Errorf("value is not an integer: 0.00"),
		},
		"key with string value": {
			envValue: "a",
			err:      fmt.Errorf("value is not an integer: a"),
		},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("1")},
			n:             1,
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
			e := &Env{
				getenv: func(key string) string {
					return tc.envValue
				},
			}
			n, err := e.Int("any", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.n, n)
		})
	}
}

func Test_GetIntRange(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		lower         int
		upper         int
		optionSetters []OptionSetter
		n             int
		err           error
	}{
		"key with int value": {
			envValue: "0",
			lower:    0,
			upper:    1,
		},
		"key with string value": {
			envValue: "a",
			lower:    0,
			upper:    1,
			err:      fmt.Errorf("value is not an integer: a"),
		},
		"key with int value below lower": {
			envValue: "0",
			lower:    1,
			upper:    2,
			err:      fmt.Errorf("value is not in range: 0 is not between 1 and 2"),
		},
		"key with int value above upper": {
			envValue: "2",
			lower:    0,
			upper:    1,
			err:      fmt.Errorf("value is not in range: 2 is not between 0 and 1"),
		},
		"key without value and default": {
			lower:         0,
			upper:         1,
			optionSetters: []OptionSetter{Default("1")},
			n:             1,
		},
		"key without value and over limit default": {
			lower:         0,
			upper:         1,
			optionSetters: []OptionSetter{Default("2")},
			err:           fmt.Errorf("value is not in range: 2 is not between 0 and 1")},
		"key without value and compulsory": {
			lower:         0,
			upper:         1,
			optionSetters: []OptionSetter{Compulsory()},
			err:           ErrNoValue,
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
			}
			n, err := e.IntRange("any", tc.lower, tc.upper, tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.n, n)
		})
	}
}
