package params

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewEnv(t *testing.T) {
	t.Parallel()
	e := NewEnv()
	assert.NotNil(t, e)
}

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
			err:           fmt.Errorf(`no value found for environment variable "key"`),
		},
		"bad options": {
			optionSetters: []OptionSetter{Compulsory(), Default("a")},
			err:           fmt.Errorf(`environment variable "key": cannot set default value for environment variable value which is compulsory`), //nolint:lll
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
			err: fmt.Errorf(`no value found for environment variable "key"`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const keyArg = "key"

			getenv := func(key string) string { return tc.env[key] }

			unsetIndex := 0
			unset := func(key string) error {
				assert.Equal(t, tc.unsetCalls[unsetIndex], key)
				unsetIndex++
				return nil
			}

			e := &envParams{
				getenv: getenv,
				unset:  unset,
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
			err:      fmt.Errorf(`environment variable "any" value "0.00" is not a valid integer`),
		},
		"key with string value": {
			envValue: "a",
			err:      fmt.Errorf(`environment variable "any" value "a" is not a valid integer`),
		},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("1")},
			n:             1,
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf(`no value found for environment variable "any"`)},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
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
			err:      fmt.Errorf(`environment variable "any" value "a" is not a valid integer`),
		},
		"key with int value below lower": {
			envValue: "0",
			lower:    1,
			upper:    2,
			err:      fmt.Errorf(`environment variable "any" value 0 is not between 1 and 2`),
		},
		"key with int value above upper": {
			envValue: "2",
			lower:    0,
			upper:    1,
			err:      fmt.Errorf(`environment variable "any" value 2 is not between 0 and 1`),
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
			err:           fmt.Errorf(`environment variable "any" value 2 is not between 0 and 1`)},
		"key without value and compulsory": {
			lower:         0,
			upper:         1,
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf(`no value found for environment variable "any"`)},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
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

func Test_YesNo(t *testing.T) { //nolint:dupl
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
			err: fmt.Errorf(`environment variable "any" value is "" and can only be "yes" or "no"`)},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("yes")},
			yes:           true,
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf(`no value found for environment variable "any"`),
		},
		"key with invalid value": {
			envValue: "a",
			err:      fmt.Errorf(`environment variable "any" value is "a" and can only be "yes" or "no"`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
				getenv: func(key string) string {
					return tc.envValue
				},
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

func Test_OnOff(t *testing.T) { //nolint:dupl
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
			err: fmt.Errorf(`environment variable "any" value is "" and can only be "on" or "off"`),
		},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("on")},
			on:            true,
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf(`no value found for environment variable "any"`),
		},
		"key with invalid value": {
			envValue: "a",
			err:      fmt.Errorf(`environment variable "any" value is "a" and can only be "on" or "off"`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
				getenv: func(key string) string {
					return tc.envValue
				},
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
			err:           fmt.Errorf(`environment variable "any" value is "a" and can only be one of: A, b`),
		},
		"key with value not in possibilities": {
			possibilities: []string{"a", "b"},
			envValue:      "c",
			err:           fmt.Errorf(`environment variable "any" value is "c" and can only be one of: a, b`),
		},
		"key without value": {
			possibilities: []string{"a", "b"},
			value:         "",
		},
		"key without value compulsory": {
			possibilities: []string{"a", "b"},
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf(`no value found for environment variable "any"`),
		},
		"key without value and default": {
			possibilities: []string{"a", "b"},
			optionSetters: []OptionSetter{Default("a")},
			value:         "a",
		},
		"key without value and compulsory": {
			possibilities: []string{"a", "b"},
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf(`no value found for environment variable "any"`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
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
			err:           fmt.Errorf(`no value found for environment variable "any"`),
		},
		"single comma": {
			envValue: ",",
			err: fmt.Errorf(
				`environment variable "any": invalid values found: value "" at position 1, value "" at position 2: possible values are: `), //nolint:lll
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
			err:           fmt.Errorf(`environment variable "any": invalid values found: value "B" at position 1: possible values are: a, b, c`), //nolint:lll
		},
		"two valid": {
			possibilities: []string{"a", "b", "c"},
			envValue:      "b,a",
			values:        []string{"b", "a"},
		},
		"one valid and one invalid": {
			possibilities: []string{"a", "b", "c"},
			envValue:      "b,d",
			err:           fmt.Errorf(`environment variable "any": invalid values found: value "d" at position 2: possible values are: a, b, c`), //nolint:lll
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
				getenv: func(key string) string {
					return tc.envValue
				},
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
			err:      fmt.Errorf(`environment variable "any" duration value is malformed: time: invalid duration "a"`),
		},
		"key without unit": {
			envValue: "1",
			err:      fmt.Errorf(`environment variable "any" duration value is malformed: time: missing unit in duration "1"`),
		},
		"key with 0 integer value": {
			envValue: "0",
			err:      fmt.Errorf(`environment variable "any" duration value cannot be 0`),
		},
		"key with negative duration": {
			envValue: "-1s",
			err:      fmt.Errorf(`environment variable "any" duration value cannot be lower than 0`),
		},
		"key without value": {
			err: fmt.Errorf(`environment variable "any" duration value is empty`),
		},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("1s")},
			duration:      time.Second,
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf(`no value found for environment variable "any"`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
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
			err:      errors.New(`no value found for environment variable "LISTENING_ADDRESS"`),
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

			e := &envParams{
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
			err:           fmt.Errorf(`port "" is not a valid integer`),
		},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("9000")},
			listeningPort: 9000,
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf(`no value found for environment variable "LISTENING_PORT"`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			const expectedUID = 1000
			e := &envParams{
				getenv: func(key string) string {
					return tc.envValue
				},
				getuid: func() int {
					return expectedUID
				},
			}
			listeningPort, warning, err := e.ListeningPort("LISTENING_PORT", tc.optionSetters...)
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

			e := &envParams{
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

func Test_RootURL(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		rootURL       string
		err           error
	}{
		"key with valid value": {
			envValue: "/a",
			rootURL:  "/a",
		},
		"key with valid value and trail /": {
			envValue: "/a/",
			rootURL:  "/a",
		},
		"key with invalid value": {
			envValue: "/a?",
			err:      fmt.Errorf(`environment variable ROOT_URL value "/a?" is not valid`),
		},
		"key without value": {},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("/a")},
			rootURL:       "/a",
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf(`environment variable "ROOT_URL": cannot make environment variable value compulsory with a default value`), //nolint:lll
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
				getenv: func(key string) string {
					return tc.envValue
				},
				regex: verification.NewRegex(),
			}
			rootURL, err := e.RootURL("ROOT_URL", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.rootURL, rootURL)
		})
	}
}

func Test_Path(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		absPath       string
		absErr        error
		path          string
		err           error
	}{
		"valid path": {
			envValue: "/path",
			absPath:  "/real/path",
			path:     "/real/path",
		},
		"get error": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           errors.New(`no value found for environment variable "key"`),
		},
		"abs error": {
			envValue: "/path",
			absErr:   errors.New("abs error"),
			err:      errors.New(`invalid filepath: for environment variable key: abs error`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
				getenv: func(key string) string {
					return tc.envValue
				},
				fpAbs: func(s string) (string, error) {
					return tc.absPath, tc.absErr
				},
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

func Test_LoggerEncoding(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		encoding      logging.Encoding
		err           error
	}{
		"key with json value": {
			envValue: "json",
			encoding: logging.JSONEncoding,
		},
		"key with console value": {
			envValue: "console",
			encoding: logging.ConsoleEncoding,
		},
		"key with invalid value": {
			envValue: "bla",
			err:      fmt.Errorf(`environment variable LOG_ENCODING: logger encoding "bla" unrecognized`),
		},
		"key without value": {
			encoding: logging.JSONEncoding,
		},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("console")},
			encoding:      logging.ConsoleEncoding,
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf(`environment variable "LOG_ENCODING": cannot make environment variable value compulsory with a default value`), //nolint:lll
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
				getenv: func(key string) string {
					return tc.envValue
				},
			}
			encoding, err := e.LoggerEncoding(tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.encoding, encoding)
		})
	}
}

func Test_LoggerLevel(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		level         logging.Level
		err           error
	}{
		"key with info value": {
			envValue: "info",
			level:    logging.InfoLevel,
		},
		"key with warning value": {
			envValue: "warning",
			level:    logging.WarnLevel,
		},
		"key with error value": {
			envValue: "error",
			level:    logging.ErrorLevel,
		},
		"key with invalid value": {
			envValue: "bla",
			level:    logging.InfoLevel,
			err:      fmt.Errorf(`environment variable LOG_LEVEL: logger level "bla" unrecognized`),
		},
		"key without value": {
			level: logging.InfoLevel,
		},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("warning")},
			level:         logging.WarnLevel,
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			level:         logging.InfoLevel,
			err:           fmt.Errorf(`environment variable "LOG_LEVEL": cannot make environment variable value compulsory with a default value`), //nolint:lll
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
				getenv: func(key string) string {
					return tc.envValue
				},
			}
			level, err := e.LoggerLevel(tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.level, level)
			}
		})
	}
}

func Test_URL(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		URL           string
		err           error
	}{
		"key with URL value": {"https://google.com", nil, "https://google.com", nil},
		// "key with invalid value":           {"please help finding me", nil, "", fmt.Errorf("")},
		"key with non HTTP value": {
			envValue: "google.com",
			err:      fmt.Errorf(`environment variable "any" URL value "google.com" is not HTTP(s)`),
		},
		"key without value": {},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("https://google.com")},
			URL:           "https://google.com",
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf(`no value found for environment variable "any"`),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &envParams{
				getenv: func(key string) string {
					return tc.envValue
				},
			}
			URL, err := e.URL("any", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			if URL == nil {
				assert.Empty(t, tc.URL)
			} else {
				assert.Equal(t, tc.URL, URL.String())
			}
		})
	}
}
