package params

import (
	"fmt"
	liburl "net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . EnvParams

// EnvParams has functions to obtain values from environment variables.
type EnvParams interface {
	GetEnv(key string, setters ...GetEnvSetter) (value string, err error)
	GetEnvInt(key string, setters ...GetEnvSetter) (n int, err error)
	GetEnvIntRange(key string, lower, upper int, setters ...GetEnvSetter) (n int, err error)
	GetYesNo(key string, setters ...GetEnvSetter) (yes bool, err error)
	GetOnOff(key string, setters ...GetEnvSetter) (on bool, err error)
	GetValueIfInside(key string, possibilities []string, setters ...GetEnvSetter) (value string, err error)
	GetCSVInPossibilities(key string, possibilities []string, setters ...GetEnvSetter) (values []string, err error)
	GetDuration(key string, setters ...GetEnvSetter) (duration time.Duration, err error)
	GetHTTPTimeout(setters ...GetEnvSetter) (duration time.Duration, err error)
	GetUserID() int
	GetGroupID() int
	GetPort(key string, setters ...GetEnvSetter) (port uint16, err error)
	GetListeningPort(key string, setters ...GetEnvSetter) (port uint16, warning string, err error)
	GetRootURL(setters ...GetEnvSetter) (rootURL string, err error)
	GetDatabaseDetails() (hostname, user, password, dbName string, err error)
	GetRedisDetails() (hostname, port, password string, err error)
	GetExeDir() (dir string, err error)
	GetPath(key string, setters ...GetEnvSetter) (path string, err error)
	GetLoggerConfig() (encoding logging.Encoding, level logging.Level, err error)
	GetLoggerEncoding(setters ...GetEnvSetter) (encoding logging.Encoding, err error)
	GetLoggerLevel(setters ...GetEnvSetter) (level logging.Level, err error)
	GetURL(key string, setters ...GetEnvSetter) (URL *liburl.URL, err error)
	GetGotifyURL(setters ...GetEnvSetter) (URL *liburl.URL, err error)
	GetGotifyToken(setters ...GetEnvSetter) (token string, err error)
}

type envParams struct {
	getenv     func(key string) string
	getuid     func() int
	getgid     func() int
	executable func() (string, error)
	verifier   verification.Verifier
	unset      func(k string) error
}

// NewEnvParams returns a new EnvParams object.
func NewEnvParams() EnvParams {
	return &envParams{
		getenv:     os.Getenv,
		getuid:     os.Getuid,
		getgid:     os.Getgid,
		executable: os.Executable,
		verifier:   verification.NewVerifier(),
		unset:      os.Unsetenv,
	}
}

type getEnvOptions struct {
	compulsory         bool
	caseSensitiveValue bool
	unset              bool
	defaultValue       string
	retroKeys          []string
	onRetro            func(oldKey, newKey string)
}

// GetEnvSetter is a setter for options to GetEnv functions.
type GetEnvSetter func(options *getEnvOptions) error

// Compulsory forces the environment variable to contain a value.
func Compulsory() GetEnvSetter {
	return func(options *getEnvOptions) error {
		if len(options.defaultValue) > 0 {
			return fmt.Errorf("cannot make environment variable value compulsory with a default value")
		}
		options.compulsory = true
		return nil
	}
}

// Default sets a default string value for the environment variable if no value is found.
func Default(defaultValue string) GetEnvSetter {
	return func(options *getEnvOptions) error {
		if options.compulsory {
			return fmt.Errorf("cannot set default value for environment variable value which is compulsory")
		}
		options.defaultValue = defaultValue
		return nil
	}
}

// CaseSensitiveValue makes the value processing case sensitive.
func CaseSensitiveValue() GetEnvSetter {
	return func(options *getEnvOptions) error {
		options.caseSensitiveValue = true
		return nil
	}
}

// Unset unsets the environment variable after it has been read.
func Unset() GetEnvSetter {
	return func(options *getEnvOptions) error {
		options.unset = true
		return nil
	}
}

// RetroKeys tries to read from retroactive environment variable keys
// and runs the function onRetro if any retro environment variable is not
// empty. RetroKeys overrides previous RetroKeys setters passed.
func RetroKeys(keys []string, onRetro func(oldKey, newKey string)) GetEnvSetter {
	return func(options *getEnvOptions) error {
		options.retroKeys = keys
		options.onRetro = onRetro
		return nil
	}
}

// GetEnv returns the value stored for a named environment variable,
// and a default if no value is found.
func (e *envParams) GetEnv(key string, setters ...GetEnvSetter) (value string, err error) {
	options := getEnvOptions{}
	defer func() {
		if options.unset {
			_ = e.unset(key)
			for _, retroKey := range options.retroKeys {
				_ = e.unset(retroKey)
			}
		}
	}()
	for _, setter := range setters {
		if err := setter(&options); err != nil {
			return "", fmt.Errorf("environment variable %q: %w", key, err)
		}
	}
	value = e.getenv(key)
	if len(value) == 0 {
		for _, retroKey := range options.retroKeys {
			value = e.getenv(retroKey)
			if len(value) > 0 {
				options.onRetro(retroKey, key)
				break
			}
		}
	}
	if len(value) == 0 {
		if options.compulsory {
			return "", fmt.Errorf("no value found for environment variable %q", key)
		}
		value = options.defaultValue
	}
	if !options.caseSensitiveValue {
		value = strings.ToLower(value)
	}
	return value, nil
}

// GetEnvInt returns the value stored for a named environment variable,
// and a default if no value is found. If the value is not a valid
// integer, an error is returned.
func (e *envParams) GetEnvInt(key string, setters ...GetEnvSetter) (n int, err error) {
	s, err := e.GetEnv(key, setters...)
	if err != nil {
		return n, err
	}
	n, err = strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("environment variable %q value %q is not a valid integer", key, s)
	}
	return n, nil
}

// GetEnvIntRange returns the value stored for a named environment variable,
// and a default if no value is found. If the value is not a valid
// integer in the range specified, an error is returned.
func (e *envParams) GetEnvIntRange(key string, lower, upper int, setters ...GetEnvSetter) (n int, err error) {
	s, err := e.GetEnv(key, setters...)
	if err != nil {
		return n, err
	}
	n, err = strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("environment variable %q value %q is not a valid integer", key, s)
	}
	if n < lower || n > upper {
		return 0, fmt.Errorf("environment variable %q value %d is not between %d and %d", key, n, lower, upper)
	}
	return n, nil
}

// GetYesNo obtains the value stored for a named environment variable and returns:
// if the value is 'yes', it returns true
// if the value is 'no', it returns false
// if it is unset, it returns the default value
// otherwise, an error is returned.
func (e *envParams) GetYesNo(key string, setters ...GetEnvSetter) (yes bool, err error) {
	s, err := e.GetEnv(key, setters...)
	if err != nil {
		return false, err
	}
	switch s {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, fmt.Errorf(`environment variable %q value is %q and can only be "yes" or "no"`, key, s)
	}
}

// GetOnOff obtains the value stored for a named environment variable and returns:
// if the value is 'on', it returns true
// if the value is 'off', it returns false
// if it is unset, it returns the default value
// otherwise, an error is returned.
func (e *envParams) GetOnOff(key string, setters ...GetEnvSetter) (on bool, err error) {
	s, err := e.GetEnv(key, setters...)
	if err != nil {
		return false, err
	}
	switch s {
	case "on":
		return true, nil
	case "off":
		return false, nil
	default:
		return false, fmt.Errorf(`environment variable %q value is %q and can only be "on" or "off"`, key, s)
	}
}

// GetValueIfInside obtains the value stored for a named environment variable if it is part of a
// list of possible values. You can optionally specify a defaultValue.
func (e *envParams) GetValueIfInside(key string, possibilities []string, setters ...GetEnvSetter) (
	value string, err error) {
	options := getEnvOptions{}
	for _, setter := range setters {
		_ = setter(&options) // error is checked in e.GetEnv
	}
	s, err := e.GetEnv(key, setters...)
	if err != nil {
		return "", err
	} else if len(s) == 0 && !options.compulsory {
		return "", nil
	}
	for _, possibility := range possibilities {
		if options.caseSensitiveValue && s == possibility {
			return s, nil
		} else if !options.caseSensitiveValue && strings.EqualFold(s, possibility) {
			return strings.ToLower(s), nil
		}
	}
	csvPossibilities := strings.Join(possibilities, ", ")
	return "", fmt.Errorf("environment variable %q value is %q and can only be one of: %s", key, s, csvPossibilities)
}

func (e *envParams) GetCSVInPossibilities(key string, possibilities []string, setters ...GetEnvSetter) (
	values []string, err error) {
	options := getEnvOptions{}
	for _, setter := range setters {
		_ = setter(&options) // error is checked in e.GetEnv
	}
	csv, err := e.GetEnv(key, setters...)
	if err != nil {
		return nil, err
	}
	if !options.compulsory && len(csv) == 0 {
		return nil, nil
	}
	values = strings.Split(csv, ",")
	type valuePosition struct {
		position int
		value    string
	}
	var invalidValues []valuePosition
	for i, value := range values {
		found := false
		for _, possibility := range possibilities {
			if options.caseSensitiveValue {
				if value == possibility {
					found = true
					break
				}
			} else {
				if strings.EqualFold(value, possibility) {
					values[i] = strings.ToLower(value)
					found = true
					break
				}
			}
		}
		if !found {
			invalidValues = append(invalidValues, valuePosition{i + 1, value})
		}
	}
	if L := len(invalidValues); L > 0 {
		invalidMessages := make([]string, L)
		for i := range invalidValues {
			invalidMessages[i] = fmt.Sprintf("value %q at position %d", invalidValues[i].value, invalidValues[i].position)
		}
		return nil, fmt.Errorf("environment variable %q: invalid values found: %s", key, strings.Join(invalidMessages, ", "))
	}
	return values, nil
}

// GetDuration gets the duration from the environment variable corresponding to the key provided.
func (e *envParams) GetDuration(key string, setters ...GetEnvSetter) (duration time.Duration, err error) {
	s, err := e.GetEnv(key, setters...)
	if err != nil {
		return 0, err
	} else if len(s) == 0 {
		return 0, fmt.Errorf("environment variable %q duration value is empty", key)
	}
	duration, err = time.ParseDuration(s)
	switch {
	case err != nil:
		return 0, fmt.Errorf("environment variable %q duration value is malformed: %w", key, err)
	case duration < 0:
		return 0, fmt.Errorf("environment variable %q duration value cannot be lower than 0", key)
	case duration == 0:
		return 0, fmt.Errorf("environment variable %q duration value cannot be 0", key)
	default:
		return duration, nil
	}
}

// GetPort obtains and checks a port number from the
// environment variable corresponding to the key provided.
func (e *envParams) GetPort(key string, setters ...GetEnvSetter) (port uint16, err error) {
	s, err := e.GetEnv(key, setters...)
	if err != nil {
		return 0, err
	}
	return verification.ParsePort(s)
}

// GetHTTPTimeout returns the HTTP client timeout duration in milliseconds
// from the environment variable HTTP_TIMEOUT.
func (e *envParams) GetHTTPTimeout(setters ...GetEnvSetter) (timeout time.Duration, err error) {
	return e.GetDuration("HTTP_TIMEOUT", setters...)
}

// GetUserID obtains the user ID running the program.
func (e *envParams) GetUserID() int {
	return e.getuid()
}

// GetGroupID obtains the user ID running the program.
func (e *envParams) GetGroupID() int {
	return e.getgid()
}

// GetListeningPort obtains and checks a port from an environment variable
// and returns a warning associated with the user ID and the port found.
func (e *envParams) GetListeningPort(key string, setters ...GetEnvSetter) (port uint16, warning string, err error) {
	port, err = e.GetPort(key, setters...)
	if err != nil {
		return 0, "", err
	}
	const (
		maxPrivilegedPort = 1023
		minDynamicPort    = 49151
	)
	if port <= maxPrivilegedPort {
		switch e.getuid() {
		case 0:
			warning = fmt.Sprintf("listening port %d allowed to be in the reserved system ports range as you are running as root", port) //nolint:lll
		case -1:
			warning = fmt.Sprintf("listening port %d allowed to be in the reserved system ports range as you are running in Windows", port) //nolint:lll
		default:
			return 0, "", fmt.Errorf("listening port %d cannot be in the reserved system ports range (1 to 1023) when running without root", port) //nolint:lll
		}
	} else if port > minDynamicPort {
		// dynamic and/or private ports.
		warning = fmt.Sprintf("listening port %d is in the dynamic/private ports range (above 49151)", port)
	}
	return port, warning, err
}

// GetRootURL obtains and checks the root URL
// from the environment variable ROOT_URL.
func (e *envParams) GetRootURL(setters ...GetEnvSetter) (rootURL string, err error) {
	setters = append([]GetEnvSetter{Default("/")}, setters...)
	rootURL, err = e.GetEnv("ROOT_URL", setters...)
	if err != nil {
		return rootURL, err
	}
	rootURL = path.Clean(rootURL)
	if !e.verifier.MatchRootURL(rootURL) {
		return "", fmt.Errorf("environment variable ROOT_URL value %q is not valid", rootURL)
	}
	rootURL = strings.TrimSuffix(rootURL, "/") // already have / from paths of router
	return rootURL, nil
}

// GetDatabaseDetails obtains the SQL database details from the
// environment variables SQL_USER, SQL_PASSWORD and SQL_DBNAME.
func (e *envParams) GetDatabaseDetails() (hostname, user, password, dbName string, err error) {
	hostname, err = e.GetEnv("SQL_HOST", Default("postgres"))
	if err != nil {
		return hostname, user, password, dbName, err
	}
	if !e.verifier.MatchHostname(hostname) {
		return hostname, user, password, dbName,
			fmt.Errorf("Postgres parameters: hostname %q is not valid", hostname)
	}
	user, err = e.GetEnv("SQL_USER", Default("postgres"), CaseSensitiveValue())
	if err != nil {
		return hostname, user, password, dbName, err
	}
	password, err = e.GetEnv("SQL_PASSWORD", Default("postgres"), CaseSensitiveValue(), Unset())
	if err != nil {
		return hostname, user, password, dbName, err
	}
	dbName, err = e.GetEnv("SQL_DBNAME", Default("postgres"), CaseSensitiveValue())
	if err != nil {
		return hostname, user, password, dbName, err
	}
	// TODO port
	return hostname, user, password, dbName, nil
}

// GetRedisDetails obtains the Redis details from the
// environment variables REDIS_HOST, REDIS_PORT and REDIS_PASSWORD.
func (e *envParams) GetRedisDetails() (hostname, port, password string, err error) {
	hostname, err = e.GetEnv("REDIS_HOST", Default("redis"))
	if err != nil {
		return hostname, port, password, err
	}
	if !e.verifier.MatchHostname(hostname) {
		return hostname, port, password,
			fmt.Errorf(`environment variable "REDIS_HOST" value %q is not valid`, hostname)
	}
	port, err = e.GetEnv("REDIS_PORT", Default("6379"))
	if err != nil {
		return hostname, port, password, err
	}
	if err := e.verifier.VerifyPort(port); err != nil {
		return hostname, port, password,
			fmt.Errorf("environment variable REDIS_PORT: %w", err)
	}
	password, err = e.GetEnv("REDIS_PASSWORD", CaseSensitiveValue(), Unset())
	if err != nil {
		return hostname, port, password, err
	}
	return hostname, port, password, nil
}

// GetExeDir obtains the executable directory.
func (e *envParams) GetExeDir() (dir string, err error) {
	ex, err := e.executable()
	if err != nil {
		return dir, err
	}
	dir = filepath.Dir(ex)
	return dir, nil
}

// GetPath obtains a path from the environment variable corresponding
// to key, and verifies it is valid. If it is a relative path,
// it prepends it with the executable path to obtain an absolute path.
// It uses defaultValue if no value is found.
func (e *envParams) GetPath(key string, setters ...GetEnvSetter) (path string, err error) {
	s, err := e.GetEnv(key, setters...)
	if err != nil {
		return "", err
	} else if !filepath.IsAbs(s) {
		exDir, err := e.GetExeDir()
		if err != nil {
			return "", err
		}
		s = filepath.Join(exDir, s)
	}
	return filepath.Abs(s)
}

// GetLoggerConfig obtains configuration details for the logger
// using the environment variables LOG_ENCODING and LOG_LEVEL.
func (e *envParams) GetLoggerConfig() (encoding logging.Encoding, level logging.Level, err error) {
	encoding, err = e.GetLoggerEncoding()
	if err != nil {
		return "", "", fmt.Errorf("logger configuration error: %w", err)
	}
	level, err = e.GetLoggerLevel()
	if err != nil {
		return "", "", fmt.Errorf("logger configuration error: %w", err)
	}
	return encoding, level, nil
}

// GetLoggerEncoding obtains the logging encoding
// from the environment variable LOG_ENCODING.
func (e *envParams) GetLoggerEncoding(setters ...GetEnvSetter) (encoding logging.Encoding, err error) {
	setters = append([]GetEnvSetter{Default("json")}, setters...)
	s, err := e.GetEnv("LOG_ENCODING", setters...)
	if err != nil {
		return "", err
	}
	s = strings.ToLower(s)
	switch s {
	case "json", "console":
		return logging.Encoding(s), nil
	default:
		return "", fmt.Errorf("environment variable LOG_ENCODING: logger encoding %q unrecognized", s)
	}
}

// GetLoggerLevel obtains the logging level
// from the environment variable LOG_LEVEL.
func (e *envParams) GetLoggerLevel(setters ...GetEnvSetter) (level logging.Level, err error) {
	setters = append([]GetEnvSetter{Default("info")}, setters...)
	s, err := e.GetEnv("LOG_LEVEL", setters...)
	if err != nil {
		return level, err
	}
	switch strings.ToLower(s) {
	case "info":
		return logging.InfoLevel, nil
	case "warning":
		return logging.WarnLevel, nil
	case "error":
		return logging.ErrorLevel, nil
	default:
		return level, fmt.Errorf("environment variable LOG_LEVEL: logger level %q unrecognized", s)
	}
}

// GetURL obtains the HTTP URL for the environment variable for the key given.
// It returns the URL of defaultValue if defaultValue is not empty.
// If no defaultValue is given, it returns nil.
func (e *envParams) GetURL(key string, setters ...GetEnvSetter) (url *liburl.URL, err error) {
	s, err := e.GetEnv(key, setters...)
	if err != nil {
		return nil, err
	} else if s == "" {
		return nil, nil
	}
	url, err = liburl.Parse(s)
	if err != nil { // never happens
		return nil, fmt.Errorf("environment variable %q URL value: %w", key, err)
	}
	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, fmt.Errorf("environment variable %q URL value %q is not HTTP(s)", key, url.String())
	}
	return url, nil
}

// GetGotifyURL obtains the URL for Gotify server
// from the environment variable GOTIFY_URL.
// It returns a nil URL if no value is found.
func (e *envParams) GetGotifyURL(setters ...GetEnvSetter) (url *liburl.URL, err error) {
	return e.GetURL("GOTIFY_URL", setters...)
}

// GetGotifyToken obtains the token for the app on the Gotify server
// from the environment variable GOTIFY_TOKEN.
func (e *envParams) GetGotifyToken(setters ...GetEnvSetter) (token string, err error) {
	setters = append(setters, CaseSensitiveValue(), Unset())
	return e.GetEnv("GOTIFY_TOKEN", setters...)
}
