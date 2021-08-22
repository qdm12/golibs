package params

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Duration(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		duration      time.Duration
		err           error
	}{
		"key with non integer value": {
			envValue: "a",
			err:      fmt.Errorf(`duration is malformed: a: time: invalid duration "a"`),
		},
		"key without unit": {
			envValue: "1",
			err:      fmt.Errorf(`duration is malformed: 1: time: missing unit in duration "1"`),
		},
		"key with 0 integer value": {
			envValue: "0",
		},
		"key with negative duration": {
			envValue: "-1s",
			err:      fmt.Errorf(`duration is negative: -1s`),
		},
		"key without value": {},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("1s")},
			duration:      time.Second,
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
			duration, err := e.Duration("any", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.duration, duration)
		})
	}
}
