package params

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Inside(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		possibilities []string
		envValue      string
		optionSetters []OptionSetter
		value         string
		err           error
	}{
		"key with value in possibilities": {
			possibilities: []string{"a", "b"},
			envValue:      "a",
			value:         "a",
		},
		"key with value in possibilities and case sensitive": {
			possibilities: []string{"a", "b"},
			envValue:      "a",
			optionSetters: []OptionSetter{CaseSensitiveValue()},
			value:         "a",
		},
		"key with value in uppercase possibilities": {
			possibilities: []string{"A", "b"},
			envValue:      "a",
			value:         "a",
		},
		"key with uppercase value in possibilities": {
			possibilities: []string{"a", "b"},
			envValue:      "A",
			value:         "a",
		},
		"key with case sensitive value in possibilities": {
			possibilities: []string{"A", "b"},
			envValue:      "a",
			optionSetters: []OptionSetter{CaseSensitiveValue()},
			err:           fmt.Errorf("value is not within the accepted values: a: it can only be one of: A, b"),
		},
		"key with value not in possibilities": {
			possibilities: []string{"a", "b"},
			envValue:      "c",
			err:           fmt.Errorf("value is not within the accepted values: c: it can only be one of: a, b"),
		},
		"key without value": {
			possibilities: []string{"a", "b"},
			value:         "",
		},
		"key without value compulsory": {
			possibilities: []string{"a", "b"},
			optionSetters: []OptionSetter{Compulsory()},
			err:           ErrNoValue,
		},
		"key without value and default": {
			possibilities: []string{"a", "b"},
			optionSetters: []OptionSetter{Default("a")},
			value:         "a",
		},
		"key without value and compulsory": {
			possibilities: []string{"a", "b"},
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
			value, err := e.Inside("any", tc.possibilities, tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.value, value)
		})
	}
}
