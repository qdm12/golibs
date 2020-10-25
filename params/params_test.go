package params

import (
	"fmt"
	"testing"
	"time"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewEnvParams(t *testing.T) {
	t.Parallel()
	e := NewEnvParams()
	assert.NotNil(t, e)
}

func Test_GetEnv(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		env     map[string]string
		setters []GetEnvSetter
		value   string
		err     error
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
			env:     map[string]string{"key": "VALUE"},
			setters: []GetEnvSetter{CaseSensitiveValue()},
			value:   "VALUE",
		},
		"key without value and default": {
			setters: []GetEnvSetter{Default("default")},
			value:   "default",
		},
		"key without value and compulsory": {
			setters: []GetEnvSetter{Compulsory()},
			err:     fmt.Errorf(`no value found for environment variable "key"`),
		},
		"bad options": {
			setters: []GetEnvSetter{Compulsory(), Default("a")},
			err:     fmt.Errorf(`environment variable "key": cannot set default value for environment variable value which is compulsory`), //nolint:lll
		},
		"retro key used": {
			env: map[string]string{
				"key":    "",
				"retro1": "",
				"retro2": "value2",
				"retro3": "value3",
			},
			setters: []GetEnvSetter{RetroKeys(
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
			setters: []GetEnvSetter{RetroKeys(
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
			setters: []GetEnvSetter{RetroKeys(
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
			e := &envParams{
				getenv: func(key string) string { return tc.env[key] },
			}
			value, err := e.GetEnv(keyArg, tc.setters...)
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

func Test_GetEnvInt(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		n        int
		err      error
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
			setters: []GetEnvSetter{Default("1")},
			n:       1,
		},
		"key without value and compulsory": {
			setters: []GetEnvSetter{Compulsory()},
			err:     fmt.Errorf(`no value found for environment variable "any"`)},
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
			n, err := e.GetEnvInt("any", tc.setters...)
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

func Test_GetEnvIntRange(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue string
		lower    int
		upper    int
		setters  []GetEnvSetter
		n        int
		err      error
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
			lower:   0,
			upper:   1,
			setters: []GetEnvSetter{Default("1")},
			n:       1,
		},
		"key without value and over limit default": {
			lower:   0,
			upper:   1,
			setters: []GetEnvSetter{Default("2")},
			err:     fmt.Errorf(`environment variable "any" value 2 is not between 0 and 1`)},
		"key without value and compulsory": {
			lower:   0,
			upper:   1,
			setters: []GetEnvSetter{Compulsory()},
			err:     fmt.Errorf(`no value found for environment variable "any"`)},
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
			n, err := e.GetEnvIntRange("any", tc.lower, tc.upper, tc.setters...)
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

func Test_GetYesNo(t *testing.T) { //nolint:dupl
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		yes      bool
		err      error
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
			setters: []GetEnvSetter{Default("yes")},
			yes:     true,
		},
		"key without value and compulsory": {
			setters: []GetEnvSetter{Compulsory()},
			err:     fmt.Errorf(`no value found for environment variable "any"`),
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
			yes, err := e.GetYesNo("any", tc.setters...)
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

func Test_GetOnOff(t *testing.T) { //nolint:dupl
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		on       bool
		err      error
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
			setters: []GetEnvSetter{Default("on")},
			on:      true,
		},
		"key without value and compulsory": {
			setters: []GetEnvSetter{Compulsory()},
			err:     fmt.Errorf(`no value found for environment variable "any"`),
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
			on, err := e.GetOnOff("any", tc.setters...)
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

func Test_GetValueIfInside(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		possibilities []string
		envValue      string
		setters       []GetEnvSetter
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
			setters:       []GetEnvSetter{CaseSensitiveValue()},
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
			setters:       []GetEnvSetter{CaseSensitiveValue()},
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
			setters:       []GetEnvSetter{Compulsory()},
			err:           fmt.Errorf(`no value found for environment variable "any"`),
		},
		"key without value and default": {
			possibilities: []string{"a", "b"},
			setters:       []GetEnvSetter{Default("a")},
			value:         "a",
		},
		"key without value and compulsory": {
			possibilities: []string{"a", "b"},
			setters:       []GetEnvSetter{Compulsory()},
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
			value, err := e.GetValueIfInside("any", tc.possibilities, tc.setters...)
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

func Test_GetCSVInPossibilities(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		possibilities []string
		envValue      string
		setters       []GetEnvSetter
		values        []string
		err           error
	}{
		"empty string": {},
		"empty string compulsory": {
			setters: []GetEnvSetter{Compulsory()},
			err:     fmt.Errorf(`no value found for environment variable "any"`),
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
			setters:       []GetEnvSetter{CaseSensitiveValue()},
			values:        []string{"B"},
		},
		"invalid case sensitive": {
			possibilities: []string{"a", "b", "c"},
			envValue:      "B",
			setters:       []GetEnvSetter{CaseSensitiveValue()},
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
			values, err := e.GetCSVInPossibilities("any", tc.possibilities, tc.setters...)
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

func Test_GetDuration(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		duration time.Duration
		err      error
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
			setters:  []GetEnvSetter{Default("1s")},
			duration: time.Second,
		},
		"key without value and compulsory": {
			setters: []GetEnvSetter{Compulsory()},
			err:     fmt.Errorf(`no value found for environment variable "any"`),
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
			duration, err := e.GetDuration("any", tc.setters...)
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

func Test_GetHTTPTimeout(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		timeout  time.Duration
		err      error
	}{
		"key with non integer value": {
			envValue: "a",
			err:      fmt.Errorf(`environment variable "HTTP_TIMEOUT" duration value is malformed: time: invalid duration "a"`),
		},
		"key without unit": {
			envValue: "1",
			err:      fmt.Errorf(`environment variable "HTTP_TIMEOUT" duration value is malformed: time: missing unit in duration "1"`), //nolint:lll
		},
		"key with 0 integer value": {
			envValue: "0",
			err:      fmt.Errorf(`environment variable "HTTP_TIMEOUT" duration value cannot be 0`),
		},
		"key with negative duration": {
			envValue: "-1s",
			err:      fmt.Errorf(`environment variable "HTTP_TIMEOUT" duration value cannot be lower than 0`),
		},
		"key without value": {
			err: fmt.Errorf(`environment variable "HTTP_TIMEOUT" duration value is empty`),
		},
		"key without value and default": {
			setters: []GetEnvSetter{Default("1s")},
			timeout: time.Second,
		},
		"key without value and compulsory": {
			envValue: "",
			setters:  []GetEnvSetter{Compulsory()},
			err:      fmt.Errorf(`no value found for environment variable "HTTP_TIMEOUT"`),
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
			timeout, err := e.GetHTTPTimeout(tc.setters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.timeout, timeout)
		})
	}
}

func Test_GetUserID(t *testing.T) {
	t.Parallel()
	const expectedUID = 1
	e := &envParams{
		getuid: func() int {
			return expectedUID
		},
	}
	uid := e.GetUserID()
	assert.Equal(t, expectedUID, uid)
}

func Test_GetGroupID(t *testing.T) {
	t.Parallel()
	const expectedUID = 1
	e := &envParams{
		getgid: func() int {
			return expectedUID
		},
	}
	gid := e.GetGroupID()
	assert.Equal(t, expectedUID, gid)
}

func Test_GetListeningPort(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		setters       []GetEnvSetter
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
			setters:       []GetEnvSetter{Default("9000")},
			listeningPort: 9000,
		},
		"key without value and compulsory": {
			setters: []GetEnvSetter{Compulsory()},
			err:     fmt.Errorf(`no value found for environment variable "LISTENING_PORT"`), //nolint:lll
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
				verifier: verification.NewVerifier(),
			}
			listeningPort, warning, err := e.GetListeningPort("LISTENING_PORT", tc.setters...)
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

func Test_GetRootURL(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		rootURL  string
		err      error
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
			setters: []GetEnvSetter{Default("/a")},
			rootURL: "/a",
		},
		"key without value and compulsory": {
			setters: []GetEnvSetter{Compulsory()},
			err:     fmt.Errorf(`environment variable "ROOT_URL": cannot make environment variable value compulsory with a default value`), //nolint:lll
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
				verifier: verification.NewVerifier(),
			}
			rootURL, err := e.GetRootURL(tc.setters...)
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

func Test_GetLoggerEncoding(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		encoding logging.Encoding
		err      error
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
			setters:  []GetEnvSetter{Default("console")},
			encoding: logging.ConsoleEncoding,
		},
		"key without value and compulsory": {
			setters: []GetEnvSetter{Compulsory()},
			err:     fmt.Errorf(`environment variable "LOG_ENCODING": cannot make environment variable value compulsory with a default value`), //nolint:lll
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
			encoding, err := e.GetLoggerEncoding(tc.setters...)
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

func Test_GetLoggerLevel(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		level    logging.Level
		err      error
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
			setters: []GetEnvSetter{Default("warning")},
			level:   logging.WarnLevel,
		},
		"key without value and compulsory": {
			setters: []GetEnvSetter{Compulsory()},
			level:   logging.InfoLevel,
			err:     fmt.Errorf(`environment variable "LOG_LEVEL": cannot make environment variable value compulsory with a default value`), //nolint:lll
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
			level, err := e.GetLoggerLevel(tc.setters...)
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

func Test_GetURL(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		URL      string
		err      error
	}{
		"key with URL value": {"https://google.com", nil, "https://google.com", nil},
		// "key with invalid value":           {"please help finding me", nil, "", fmt.Errorf("")},
		"key with non HTTP value": {
			envValue: "google.com",
			err:      fmt.Errorf(`environment variable "any" URL value "google.com" is not HTTP(s)`),
		},
		"key without value": {},
		"key without value and default": {
			setters: []GetEnvSetter{Default("https://google.com")},
			URL:     "https://google.com",
		},
		"key without value and compulsory": {
			setters: []GetEnvSetter{Compulsory()},
			err:     fmt.Errorf(`no value found for environment variable "any"`),
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
			URL, err := e.GetURL("any", tc.setters...)
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

func Test_GetGotifyURLL(t *testing.T) {
	t.Parallel()
	e := &envParams{
		getenv: func(key string) string {
			return "https://google.com"
		},
	}
	URL, err := e.GetGotifyURL()
	require.NoError(t, err)
	require.NotNil(t, URL)
	assert.Equal(t, "https://google.com", URL.String())
}

func Test_GetGotifyToken(t *testing.T) {
	t.Parallel()
	e := &envParams{
		getenv: func(key string) string {
			return "x"
		},
		unset: func(k string) error { return nil },
	}
	token, err := e.GetGotifyToken()
	require.NoError(t, err)
	assert.Equal(t, "x", token)
}
