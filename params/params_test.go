package params

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
)

func Test_NewEnvParams(t *testing.T) {
	t.Parallel()
	e := NewEnvParams()
	assert.NotNil(t, e)
}

func Test_GetEnv(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		value    string
		err      error
	}{
		"key with value":                   {"value", nil, "value", nil},
		"key with uppercase value":         {"VALUE", nil, "value", nil},
		"key with case sensitive value":    {"VALUE", []GetEnvSetter{CaseSensitiveValue()}, "VALUE", nil},
		"key without value and default":    {"", []GetEnvSetter{Default("default")}, "default", nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, "", fmt.Errorf("no value found for environment variable \"any\"")},
		"bad options":                      {"", []GetEnvSetter{Compulsory(), Default("a")}, "", fmt.Errorf("environment variable \"any\": cannot set default value for environment variable value which is compulsory")},
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
			value, err := e.GetEnv("any", tc.setters...)
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
		"key with int value":               {"0", nil, 0, nil},
		"key with float value":             {"0.00", nil, 0, fmt.Errorf("environment variable \"any\" value \"0.00\" is not a valid integer")},
		"key with string value":            {"a", nil, 0, fmt.Errorf("environment variable \"any\" value \"a\" is not a valid integer")},
		"key without value and default":    {"", []GetEnvSetter{Default("1")}, 1, nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, 0, fmt.Errorf("no value found for environment variable \"any\"")},
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
		"key with int value":                       {"0", 0, 1, nil, 0, nil},
		"key with string value":                    {"a", 0, 1, nil, 0, fmt.Errorf("environment variable \"any\" value \"a\" is not a valid integer")},
		"key with int value below lower":           {"0", 1, 2, nil, 0, fmt.Errorf("environment variable \"any\" value 0 is not between 1 and 2")},
		"key with int value above upper":           {"2", 0, 1, nil, 0, fmt.Errorf("environment variable \"any\" value 2 is not between 0 and 1")},
		"key without value and default":            {"", 0, 1, []GetEnvSetter{Default("1")}, 1, nil},
		"key without value and over limit default": {"", 0, 1, []GetEnvSetter{Default("2")}, 0, fmt.Errorf("environment variable \"any\" value 2 is not between 0 and 1")},
		"key without value and compulsory":         {"", 0, 1, []GetEnvSetter{Compulsory()}, 0, fmt.Errorf("no value found for environment variable \"any\"")},
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
		"key with yes value":               {"yes", nil, true, nil},
		"key with no value":                {"no", nil, false, nil},
		"key without value":                {"", nil, false, fmt.Errorf("environment variable \"any\" value is \"\" and can only be \"yes\" or \"no\"")},
		"key without value and default":    {"", []GetEnvSetter{Default("yes")}, true, nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, false, fmt.Errorf("no value found for environment variable \"any\"")},
		"key with invalid value":           {"a", nil, false, fmt.Errorf("environment variable \"any\" value is \"a\" and can only be \"yes\" or \"no\"")},
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
		"key with on value":                {"on", nil, true, nil},
		"key with off value":               {"off", nil, false, nil},
		"key without value":                {"", nil, false, fmt.Errorf("environment variable \"any\" value is \"\" and can only be \"on\" or \"off\"")},
		"key without value and default":    {"", []GetEnvSetter{Default("on")}, true, nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, false, fmt.Errorf("no value found for environment variable \"any\"")},
		"key with invalid value":           {"a", nil, false, fmt.Errorf("environment variable \"any\" value is \"a\" and can only be \"on\" or \"off\"")},
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
		"key with value in possibilities":                {[]string{"a", "b"}, "a", nil, "a", nil},
		"key with value in uppercase possibilities":      {[]string{"A", "b"}, "a", nil, "a", nil},
		"key with uppercase value in possibilities":      {[]string{"a", "b"}, "A", nil, "a", nil},
		"key with case sensitive value in possibilities": {[]string{"A", "b"}, "a", []GetEnvSetter{CaseSensitiveValue()}, "", fmt.Errorf("environment variable \"any\" value is \"a\" and can only be one of: A, b")},
		"key with value not in possibilities":            {[]string{"a", "b"}, "c", nil, "", fmt.Errorf("environment variable \"any\" value is \"c\" and can only be one of: a, b")},
		"key without value":                              {[]string{"a", "b"}, "", nil, "", fmt.Errorf("environment variable \"any\" value is \"\" and can only be one of: a, b")},
		"key without value and default":                  {[]string{"a", "b"}, "", []GetEnvSetter{Default("a")}, "a", nil},
		"key without value and compulsory":               {[]string{"a", "b"}, "", []GetEnvSetter{Compulsory()}, "", fmt.Errorf("no value found for environment variable \"any\"")},
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

func Test_GetDuration(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		duration time.Duration
		err      error
	}{
		"key with non integer value":       {"a", nil, 0, fmt.Errorf("environment variable \"any\" duration value is malformed: time: invalid duration a")},
		"key without unit":                 {"1", nil, 0, fmt.Errorf("environment variable \"any\" duration value is malformed: time: missing unit in duration 1")},
		"key with 0 integer value":         {"0", nil, 0, fmt.Errorf("environment variable \"any\" duration value cannot be 0")},
		"key with negative duration":       {"-1s", nil, 0, fmt.Errorf("environment variable \"any\" duration value cannot be lower than 0")},
		"key without value":                {"", nil, 0, fmt.Errorf("environment variable \"any\" duration value is empty")},
		"key without value and default":    {"", []GetEnvSetter{Default("1s")}, time.Second, nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, 0, fmt.Errorf("no value found for environment variable \"any\"")},
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
		"key with non integer value":       {"a", nil, 0, fmt.Errorf("environment variable \"HTTP_TIMEOUT\" duration value is malformed: time: invalid duration a")},
		"key without unit":                 {"1", nil, 0, fmt.Errorf("environment variable \"HTTP_TIMEOUT\" duration value is malformed: time: missing unit in duration 1")},
		"key with 0 integer value":         {"0", nil, 0, fmt.Errorf("environment variable \"HTTP_TIMEOUT\" duration value cannot be 0")},
		"key with negative duration":       {"-1s", nil, 0, fmt.Errorf("environment variable \"HTTP_TIMEOUT\" duration value cannot be lower than 0")},
		"key without value":                {"", nil, 0, fmt.Errorf("environment variable \"HTTP_TIMEOUT\" duration value is empty")},
		"key without value and default":    {"", []GetEnvSetter{Default("1s")}, time.Second, nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, 0, fmt.Errorf("no value found for environment variable \"HTTP_TIMEOUT\"")},
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
		listeningPort string
		warning       string
		err           error
	}{
		"key with valid value":             {"9000", nil, "9000", "", nil},
		"key with valid warning value":     {"60000", nil, "60000", "listening port 60000 is in the dynamic/private ports range (above 49151)", nil},
		"key without value":                {"", nil, "8000", "", nil},
		"key without value and default":    {"", []GetEnvSetter{Default("9000")}, "9000", "", nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, "", "", fmt.Errorf("environment variable \"LISTENING_PORT\": cannot make environment variable value compulsory with a default value")},
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
			listeningPort, warning, err := e.GetListeningPort(tc.setters...)
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
		"key with valid value":             {"/a", nil, "/a", nil},
		"key with valid value and trail /": {"/a/", nil, "/a", nil},
		"key with invalid value":           {"/a?", nil, "", fmt.Errorf("environment variable ROOT_URL value \"/a?\" is not valid")},
		"key without value":                {"", nil, "", nil},
		"key without value and default":    {"", []GetEnvSetter{Default("/a")}, "/a", nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, "", fmt.Errorf("environment variable \"ROOT_URL\": cannot make environment variable value compulsory with a default value")},
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
		"key with json value":              {"json", nil, logging.JSONEncoding, nil},
		"key with console value":           {"console", nil, logging.ConsoleEncoding, nil},
		"key with invalid value":           {"bla", nil, "", fmt.Errorf("environment variable LOG_ENCODING: logger encoding \"bla\" unrecognized")},
		"key without value":                {"", nil, logging.JSONEncoding, nil},
		"key without value and default":    {"", []GetEnvSetter{Default("console")}, logging.ConsoleEncoding, nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, "", fmt.Errorf("environment variable \"LOG_ENCODING\": cannot make environment variable value compulsory with a default value")},
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
		"key with info value":              {"info", nil, logging.InfoLevel, nil},
		"key with warning value":           {"warning", nil, logging.WarnLevel, nil},
		"key with error value":             {"error", nil, logging.ErrorLevel, nil},
		"key with invalid value":           {"bla", nil, logging.InfoLevel, fmt.Errorf("environment variable LOG_LEVEL: logger level \"bla\" unrecognized")},
		"key without value":                {"", nil, logging.InfoLevel, nil},
		"key without value and default":    {"", []GetEnvSetter{Default("warning")}, logging.WarnLevel, nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, logging.InfoLevel, fmt.Errorf("environment variable \"LOG_LEVEL\": cannot make environment variable value compulsory with a default value")},
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

func Test_GetNodeID(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue string
		setters  []GetEnvSetter
		nodeID   int
		err      error
	}{
		"key with value 10":                {"10", nil, 10, nil},
		"key with invalid value":           {"bla", nil, 0, fmt.Errorf("environment variable NODE_ID value \"bla\" is not a valid integer")},
		"key without value":                {"", nil, 0, nil},
		"key without value and default":    {"", []GetEnvSetter{Default("2")}, 2, nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, 0, fmt.Errorf("environment variable \"NODE_ID\": cannot make environment variable value compulsory with a default value")},
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
			nodeID, err := e.GetNodeID(tc.setters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.nodeID, nodeID)
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
		"key with non HTTP value":          {"google.com", nil, "", fmt.Errorf("environment variable \"any\" URL value \"google.com\" is not HTTP(s)")},
		"key without value":                {"", nil, "", nil},
		"key without value and default":    {"", []GetEnvSetter{Default("https://google.com")}, "https://google.com", nil},
		"key without value and compulsory": {"", []GetEnvSetter{Compulsory()}, "", fmt.Errorf("no value found for environment variable \"any\"")},
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
