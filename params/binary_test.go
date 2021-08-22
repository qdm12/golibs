package params

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_YesNo(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		yes           bool
		err           error
	}{
		"key with yes value": {
			envValue: "yes",
			yes:      true,
		},
		"key with no value": {
			envValue: "no",
		},
		"key without value": {
			err: fmt.Errorf(`value can only be 'yes' or 'no': `)},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("yes")},
			yes:           true,
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           ErrNoValue,
		},
		"key with invalid value": {
			envValue: "a",
			err:      fmt.Errorf(`value can only be 'yes' or 'no': a`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &Env{
				kv: map[string]string{"any": tc.envValue},
			}
			yes, err := e.YesNo("any", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.yes, yes)
		})
	}
}

func Test_OnOff(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		on            bool
		err           error
	}{
		"key with on value": {
			envValue: "on", on: true,
		},
		"key with off value": {
			envValue: "off",
		},
		"key without value": {
			err: fmt.Errorf(`value can only be 'on' or 'off': `),
		},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("on")},
			on:            true,
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           ErrNoValue,
		},
		"key with invalid value": {
			envValue: "a",
			err:      fmt.Errorf(`value can only be 'on' or 'off': a`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &Env{
				kv: map[string]string{"any": tc.envValue},
			}
			on, err := e.OnOff("any", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.on, on)
		})
	}
}
