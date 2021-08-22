package params

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Get(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		env           map[string]string
		optionSetters []OptionSetter
		unsetCalls    []string
		value         string
		err           error
	}{
		"key with value": {
			env:   map[string]string{"key": "value"},
			value: "value",
		},
		"key with uppercase value": {
			env:   map[string]string{"key": "VALUE"},
			value: "value",
		},
		"key with case sensitive value": {
			env:           map[string]string{"key": "VALUE"},
			optionSetters: []OptionSetter{CaseSensitiveValue()},
			value:         "VALUE",
		},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("default")},
			value:         "default",
		},
		"key without value and unset": {
			env: map[string]string{"key": "VALUE"},
			optionSetters: []OptionSetter{
				Unset(),
				RetroKeys(
					[]string{"retro"},
					func(oldKey string, newKey string) {},
				),
			},
			unsetCalls: []string{"key", "retro"},
			value:      "value",
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           ErrNoValue,
		},
		"bad options": {
			optionSetters: []OptionSetter{Compulsory(), Default("a")},
			err:           fmt.Errorf("option error: cannot set a default for a compulsory environment variable value"),
		},
		"retro key used": {
			env: map[string]string{
				"key":    "",
				"retro1": "",
				"retro2": "value2",
				"retro3": "value3",
			},
			optionSetters: []OptionSetter{RetroKeys(
				[]string{"retro1", "retro2", "retro3"},
				func(oldKey string, newKey string) {
					assert.Equal(t, "retro2", oldKey)
					assert.Equal(t, "key", newKey)
				},
			)},
			value: "value2",
		},
		"retro key unused": {
			env: map[string]string{
				"key":    "value",
				"retro1": "value1",
			},
			optionSetters: []OptionSetter{RetroKeys(
				[]string{"retro1"},
				func(oldKey string, newKey string) {},
			)},
			value: "value",
		},
		"not found with retro key": {
			env: map[string]string{
				"key":    "",
				"retro1": "",
			},
			optionSetters: []OptionSetter{RetroKeys(
				[]string{"retro1"},
				func(oldKey string, newKey string) {},
			), Compulsory()},
			err: ErrNoValue,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const keyArg = "key"

			unsetIndex := 0
			unset := func(key string) error {
				assert.Equal(t, tc.unsetCalls[unsetIndex], key)
				unsetIndex++
				return nil
			}

			e := &Env{
				kv:    tc.env,
				unset: unset,
			}
			value, err := e.Get(keyArg, tc.optionSetters...)
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

func Test_getEnv(t *testing.T) {
	t.Parallel()

	const someEnvKey = "TESTREMOVEME32932"

	env := &Env{
		kv: map[string]string{"a": "value1"},
	}

	value := env.getEnv("a")
	assert.Equal(t, "value1", value)

	err := os.Setenv(someEnvKey, "value2")
	require.NoError(t, err)
	defer func() {
		err := os.Unsetenv(someEnvKey)
		assert.NoError(t, err)
	}()

	value = env.getEnv(someEnvKey)
	assert.Equal(t, "value2", value)
}
