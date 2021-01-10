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

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Env

// Env has functions to obtain values from environment variables.
type Env interface {
	Get(key string, optionSetters ...OptionSetter) (value string, err error)
	Int(key string, optionSetters ...OptionSetter) (n int, err error)
	IntRange(key string, lower, upper int, optionSetters ...OptionSetter) (n int, err error)
	YesNo(key string, optionSetters ...OptionSetter) (yes bool, err error)
	OnOff(key string, optionSetters ...OptionSetter) (on bool, err error)
	Inside(key string, possibilities []string, optionSetters ...OptionSetter) (value string, err error)
	CSVInside(key string, possibilities []string, optionSetters ...OptionSetter) (values []string, err error)
	Duration(key string, optionSetters ...OptionSetter) (duration time.Duration, err error)
	HTTPTimeout(optionSetters ...OptionSetter) (duration time.Duration, err error)
	UserID() int
	GroupID() int
	Port(key string, optionSetters ...OptionSetter) (port uint16, err error)
	ListeningPort(key string, optionSetters ...OptionSetter) (port uint16, warning string, err error)
	RootURL(optionSetters ...OptionSetter) (rootURL string, err error)
	DatabaseDetails() (hostname, user, password, dbName string, err error)
	RedisDetails() (hostname, port, password string, err error)
	ExeDir() (dir string, err error)
	Path(key string, optionSetters ...OptionSetter) (path string, err error)
	LoggerConfig() (encoding logging.Encoding, level logging.Level, err error)
	LoggerEncoding(optionSetters ...OptionSetter) (encoding logging.Encoding, err error)
	LoggerLevel(optionSetters ...OptionSetter) (level logging.Level, err error)
	URL(key string, optionSetters ...OptionSetter) (URL *liburl.URL, err error)
	GotifyURL(optionSetters ...OptionSetter) (URL *liburl.URL, err error)
	GotifyToken(optionSetters ...OptionSetter) (token string, err error)
}

type envParams struct {
	getenv     func(key string) string
	getuid     func() int
	getgid     func() int
	executable func() (string, error)
	verifier   verification.Verifier
	unset      func(k string) error
}

// NewEnv returns a new Env object.
func NewEnv() Env {
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

// OptionSetter is a setter for options to Get functions.
type OptionSetter func(options *getEnvOptions) error

// Compulsory forces the environment variable to contain a value.
func Compulsory() OptionSetter {
	return func(options *getEnvOptions) error {
		if len(options.defaultValue) > 0 {
			return fmt.Errorf("cannot make environment variable value compulsory with a default value")
		}
		options.compulsory = true
		return nil
	}
}

// Default sets a default string value for the environment variable if no value is found.
func Default(defaultValue string) OptionSetter {
	return func(options *getEnvOptions) error {
		if options.compulsory {
			return fmt.Errorf("cannot set default value for environment variable value which is compulsory")
		}
		options.defaultValue = defaultValue
		return nil
	}
}

// CaseSensitiveValue makes the value processing case sensitive.
func CaseSensitiveValue() OptionSetter {
	return func(options *getEnvOptions) error {
		options.caseSensitiveValue = true
		return nil
	}
}

// Unset unsets the environment variable after it has been read.
func Unset() OptionSetter {
	return func(options *getEnvOptions) error {
		options.unset = true
		return nil
	}
}

// RetroKeys tries to read from retroactive environment variable keys
// and runs the function onRetro if any retro environment variable is not
// empty. RetroKeys overrides previous RetroKeys optionSetters passed.
func RetroKeys(keys []string, onRetro func(oldKey, newKey string)) OptionSetter {
	return func(options *getEnvOptions) error {
		options.retroKeys = keys
		options.onRetro = onRetro
		return nil
	}
}

// Get returns the value stored for a named environment variable,
// and a default if no value is found.
func (e *envParams) Get(key string, optionSetters ...OptionSetter) (value string, err error) {
	options := getEnvOptions{}
	defer func() {
		if options.unset {
			_ = e.unset(key)
			for _, retroKey := range options.retroKeys {
				_ = e.unset(retroKey)
			}
		}
	}()
	for _, setter := range optionSetters {
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

// Int returns the value stored for a named environment variable,
// and a default if no value is found. If the value is not a valid
// integer, an error is returned.
func (e *envParams) Int(key string, optionSetters ...OptionSetter) (n int, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return n, err
	}
	n, err = strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("environment variable %q value %q is not a valid integer", key, s)
	}
	return n, nil
}

// IntRange returns the value stored for a named environment variable,
// and a default if no value is found. If the value is not a valid
// integer in the range specified, an error is returned.
func (e *envParams) IntRange(key string, lower, upper int, optionSetters ...OptionSetter) (n int, err error) {
	s, err := e.Get(key, optionSetters...)
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

// YesNo obtains the value stored for a named environment variable and returns:
// if the value is 'yes', it returns true
// if the value is 'no', it returns false
// if it is unset, it returns the default value
// otherwise, an error is returned.
func (e *envParams) YesNo(key string, optionSetters ...OptionSetter) (yes bool, err error) {
	s, err := e.Get(key, optionSetters...)
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

// OnOff obtains the value stored for a named environment variable and returns:
// if the value is 'on', it returns true
// if the value is 'off', it returns false
// if it is unset, it returns the default value
// otherwise, an error is returned.
func (e *envParams) OnOff(key string, optionSetters ...OptionSetter) (on bool, err error) {
	s, err := e.Get(key, optionSetters...)
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

// Inside obtains the value stored for a named environment variable if it is part of a
// list of possible values. You can optionally specify a defaultValue.
func (e *envParams) Inside(key string, possibilities []string, optionSetters ...OptionSetter) (
	value string, err error) {
	options := getEnvOptions{}
	for _, setter := range optionSetters {
		_ = setter(&options) // error is checked in e.Get
	}
	s, err := e.Get(key, optionSetters...)
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

func (e *envParams) CSVInside(key string, possibilities []string, optionSetters ...OptionSetter) (
	values []string, err error) {
	options := getEnvOptions{}
	for _, setter := range optionSetters {
		_ = setter(&options) // error is checked in e.Get
	}
	csv, err := e.Get(key, optionSetters...)
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
		csvPossibilities := strings.Join(possibilities, ", ")
		return nil, fmt.Errorf("environment variable %q: invalid values found: %s: possible values are: %s",
			key, strings.Join(invalidMessages, ", "), csvPossibilities)
	}
	return values, nil
}

// Duration gets the duration from the environment variable corresponding to the key provided.
func (e *envParams) Duration(key string, optionSetters ...OptionSetter) (duration time.Duration, err error) {
	s, err := e.Get(key, optionSetters...)
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

// Port obtains and checks a port number from the
// environment variable corresponding to the key provided.
func (e *envParams) Port(key string, optionSetters ...OptionSetter) (port uint16, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return 0, err
	}
	return verification.ParsePort(s)
}

// HTTPTimeout returns the HTTP client timeout duration in milliseconds
// from the environment variable HTTP_TIMEOUT.
func (e *envParams) HTTPTimeout(optionSetters ...OptionSetter) (timeout time.Duration, err error) {
	return e.Duration("HTTP_TIMEOUT", optionSetters...)
}

// UserID obtains the user ID running the program.
func (e *envParams) UserID() int {
	return e.getuid()
}

// GroupID obtains the user ID running the program.
func (e *envParams) GroupID() int {
	return e.getgid()
}

// ListeningPort obtains and checks a port from an environment variable
// and returns a warning associated with the user ID and the port found.
func (e *envParams) ListeningPort(key string, optionSetters ...OptionSetter) (port uint16, warning string, err error) {
	port, err = e.Port(key, optionSetters...)
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

// RootURL obtains and checks the root URL from the environment variable specified by envKey.
func (e *envParams) RootURL(optionSetters ...OptionSetter) (rootURL string, err error) {
	optionSetters = append([]OptionSetter{Default("/")}, optionSetters...)
	rootURL, err = e.Get("ROOT_URL", optionSetters...)
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

// DatabaseDetails obtains the SQL database details from the
// environment variables SQL_USER, SQL_PASSWORD and SQL_DBNAME.
func (e *envParams) DatabaseDetails() (hostname, user, password, dbName string, err error) {
	hostname, err = e.Get("SQL_HOST", Default("postgres"))
	if err != nil {
		return hostname, user, password, dbName, err
	}
	if !e.verifier.MatchHostname(hostname) {
		return hostname, user, password, dbName,
			fmt.Errorf("Postgres parameters: hostname %q is not valid", hostname)
	}
	user, err = e.Get("SQL_USER", Default("postgres"), CaseSensitiveValue())
	if err != nil {
		return hostname, user, password, dbName, err
	}
	password, err = e.Get("SQL_PASSWORD", Default("postgres"), CaseSensitiveValue(), Unset())
	if err != nil {
		return hostname, user, password, dbName, err
	}
	dbName, err = e.Get("SQL_DBNAME", Default("postgres"), CaseSensitiveValue())
	if err != nil {
		return hostname, user, password, dbName, err
	}
	// TODO port
	return hostname, user, password, dbName, nil
}

// RedisDetails obtains the Redis details from the
// environment variables REDIS_HOST, REDIS_PORT and REDIS_PASSWORD.
func (e *envParams) RedisDetails() (hostname, port, password string, err error) {
	hostname, err = e.Get("REDIS_HOST", Default("redis"))
	if err != nil {
		return hostname, port, password, err
	}
	if !e.verifier.MatchHostname(hostname) {
		return hostname, port, password,
			fmt.Errorf(`environment variable "REDIS_HOST" value %q is not valid`, hostname)
	}
	port, err = e.Get("REDIS_PORT", Default("6379"))
	if err != nil {
		return hostname, port, password, err
	}
	if err := e.verifier.VerifyPort(port); err != nil {
		return hostname, port, password,
			fmt.Errorf("environment variable REDIS_PORT: %w", err)
	}
	password, err = e.Get("REDIS_PASSWORD", CaseSensitiveValue(), Unset())
	if err != nil {
		return hostname, port, password, err
	}
	return hostname, port, password, nil
}

// ExeDir obtains the executable directory.
func (e *envParams) ExeDir() (dir string, err error) {
	ex, err := e.executable()
	if err != nil {
		return dir, err
	}
	dir = filepath.Dir(ex)
	return dir, nil
}

// Path obtains a path from the environment variable corresponding
// to key, and verifies it is valid. If it is a relative path,
// it prepends it with the executable path to obtain an absolute path.
// It uses defaultValue if no value is found.
func (e *envParams) Path(key string, optionSetters ...OptionSetter) (path string, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return "", err
	} else if !filepath.IsAbs(s) {
		exDir, err := e.ExeDir()
		if err != nil {
			return "", err
		}
		s = filepath.Join(exDir, s)
	}
	return filepath.Abs(s)
}

// LoggerConfig obtains configuration details for the logger
// using the environment variables LOG_ENCODING and LOG_LEVEL.
func (e *envParams) LoggerConfig() (encoding logging.Encoding, level logging.Level, err error) {
	encoding, err = e.LoggerEncoding()
	if err != nil {
		return "", "", fmt.Errorf("logger configuration error: %w", err)
	}
	level, err = e.LoggerLevel()
	if err != nil {
		return "", "", fmt.Errorf("logger configuration error: %w", err)
	}
	return encoding, level, nil
}

// LoggerEncoding obtains the logging encoding
// from the environment variable LOG_ENCODING.
func (e *envParams) LoggerEncoding(optionSetters ...OptionSetter) (encoding logging.Encoding, err error) {
	optionSetters = append([]OptionSetter{Default("json")}, optionSetters...)
	s, err := e.Get("LOG_ENCODING", optionSetters...)
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

// LoggerLevel obtains the logging level
// from the environment variable LOG_LEVEL.
func (e *envParams) LoggerLevel(optionSetters ...OptionSetter) (level logging.Level, err error) {
	optionSetters = append([]OptionSetter{Default("info")}, optionSetters...)
	s, err := e.Get("LOG_LEVEL", optionSetters...)
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

// URL obtains the HTTP URL for the environment variable for the key given.
// It returns the URL of defaultValue if defaultValue is not empty.
// If no defaultValue is given, it returns nil.
func (e *envParams) URL(key string, optionSetters ...OptionSetter) (url *liburl.URL, err error) {
	s, err := e.Get(key, optionSetters...)
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

// GotifyURL obtains the URL for Gotify server
// from the environment variable GOTIFY_URL.
// It returns a nil URL if no value is found.
func (e *envParams) GotifyURL(optionSetters ...OptionSetter) (url *liburl.URL, err error) {
	return e.URL("GOTIFY_URL", optionSetters...)
}

// GotifyToken obtains the token for the app on the Gotify server
// from the environment variable GOTIFY_TOKEN.
func (e *envParams) GotifyToken(optionSetters ...OptionSetter) (token string, err error) {
	optionSetters = append(optionSetters, CaseSensitiveValue(), Unset())
	return e.Get("GOTIFY_TOKEN", optionSetters...)
}
