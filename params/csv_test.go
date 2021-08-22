package params

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CSVInside(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		possibilities []string
		envValue      string
		optionSetters []OptionSetter
		values        []string
		err           error
	}{
		"empty string": {},
		"empty string compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           ErrNoValue,
		},
		"single comma": {
			envValue: ",",
			err: fmt.Errorf(
				`at least one value is not within the accepted values: invalid values found: value "" at position 1, value "" at position 2; possible values are: `), //nolint:lll
		},
		"single valid": {
			possibilities: []string{"a", "b", "c"},
			envValue:      "B",
			values:        []string{"b"},
		},
		"single valid case sensitive": {
			possibilities: []string{"a", "B", "c"},
			envValue:      "B",
			optionSetters: []OptionSetter{CaseSensitiveValue()},
			values:        []string{"B"},
		},
		"invalid case sensitive": {
			possibilities: []string{"a", "b", "c"},
			envValue:      "B",
			optionSetters: []OptionSetter{CaseSensitiveValue()},
			err:           fmt.Errorf(`at least one value is not within the accepted values: invalid values found: value "B" at position 1; possible values are: a, b, c`), //nolint:lll
		},
		"two valid": {
			possibilities: []string{"a", "b", "c"},
			envValue:      "b,a",
			values:        []string{"b", "a"},
		},
		"one valid and one invalid": {
			possibilities: []string{"a", "b", "c"},
			envValue:      "b,d",
			err:           fmt.Errorf(`at least one value is not within the accepted values: invalid values found: value "d" at position 2; possible values are: a, b, c`), //nolint:lll
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &Env{
				kv: map[string]string{"any": tc.envValue},
			}
			values, err := e.CSVInside("any", tc.possibilities, tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.values, values)
		})
	}
}
