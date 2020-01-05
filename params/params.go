package params

import (
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"fmt"

	"github.com/qdm12/golibs/logging"
	"github.com/qdm12/golibs/verification"
)

// EnvParams has functions to obtain values from environment variables
type EnvParams interface {
	GetEnv(key string, setters ...GetEnvSetter) (value string, err error)
	GetEnvInt(key string, setters ...GetEnvSetter) (n int, err error)
	GetEnvIntRange(key string, lower, upper int, setters ...GetEnvSetter) (n int, err error)
	GetYesNo(key string, setters ...GetEnvSetter) (yes bool, err error)
	GetOnOff(key string, setters ...GetEnvSetter) (on bool, err error)
	GetValueIfInside(key string, possibilities []string, setters ...GetEnvSetter) (value string, err error)
	GetDuration(key string, timeUnit time.Duration, setters ...GetEnvSetter) (duration time.Duration, err error)
	GetHTTPTimeout(timeUnit time.Duration, setters ...GetEnvSetter) (duration time.Duration, err error)
	GetUserID() int
	GetGroupID() int
	GetListeningPort(setters ...GetEnvSetter) (listeningPort string, err error)
	GetRootURL(setters ...GetEnvSetter) (rootURL string, err error)
	GetDatabaseDetails() (hostname, user, password, dbName string, err error)
	GetRedisDetails() (hostname, port, password string, err error)
	GetExeDir() (dir string, err error)
	GetPath(key string, setters ...GetEnvSetter) (path string, err error)
	GetLoggerConfig() (encoding string, level logging.Level, nodeID int, err error)
	GetLoggerEncoding(setters ...GetEnvSetter) (encoding string, err error)
	GetLoggerLevel(setters ...GetEnvSetter) (level logging.Level, err error)
	GetNodeID(setters ...GetEnvSetter) (nodeID int, err error)
	GetURL(key string, setters ...GetEnvSetter) (URL *url.URL, err error)
	GetGotifyURL(setters ...GetEnvSetter) (URL *url.URL, err error)
	GetGotifyToken(setters ...GetEnvSetter) (token string, err error)
}

type envParamsImpl struct {
	getenv     func(key string) string
	getuid     func() int
	getgid     func() int
	executable func() (string, error)
}

// NewEnvParams returns a new EnvParams object
func NewEnvParams() EnvParams {
	return &envParamsImpl{
		getenv:     os.Getenv,
		getuid:     os.Getuid,
		getgid:     os.Getgid,
		executable: os.Executable,
	}
}

type getEnvOptions struct {
	compulsory   bool
	defaultValue string
}

// GetEnvSetter is a setter for options to GetEnv functions
type GetEnvSetter func(options *getEnvOptions) error

// Compulsory forces the environment variable to contain a value
func Compulsory() GetEnvSetter {
	return func(options *getEnvOptions) error {
		if len(options.defaultValue) > 0 {
			return fmt.Errorf("cannot make environment variable value compulsory with a default value")
		}
		options.compulsory = true
		return nil
	}
}

// Default sets a default string value for the environment variable if no value is found
func Default(defaultValue string) GetEnvSetter {
	return func(options *getEnvOptions) error {
		if options.compulsory {
			return fmt.Errorf("cannot set default value for environment variable value which is compulsory")
		}
		options.defaultValue = defaultValue
		return nil
	}
}

// GetEnv returns the value stored for a named environment variable,
// and a default if no value is found
func (e *envParamsImpl) GetEnv(key string, setters ...GetEnvSetter) (value string, err error) {
	options := getEnvOptions{}
	for _, setter := range setters {
		if err := setter(&options); err != nil {
			return "", fmt.Errorf("environment variable %q: %w", key, err)
		}
	}
	value = e.getenv(key)
	if len(value) == 0 {
		if options.compulsory {
			return "", fmt.Errorf("no value found for environment variable %q", key)
		}
		value = options.defaultValue
	}
	return value, nil
}

// GetEnvInt returns the value stored for a named environment variable,
// and a default if no value is found. If the value is not a valid
// integer, an error is returned.
func (e *envParamsImpl) GetEnvInt(key string, setters ...GetEnvSetter) (n int, err error) {
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
func (e *envParamsImpl) GetEnvIntRange(key string, lower, upper int, setters ...GetEnvSetter) (n int, err error) {
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
func (e *envParamsImpl) GetYesNo(key string, setters ...GetEnvSetter) (yes bool, err error) {
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
		return false, fmt.Errorf("environment variable %q value is %q and can only be \"yes\" or \"no\"", key, s)
	}
}

// GetOnOff obtains the value stored for a named environment variable and returns:
// if the value is 'on', it returns true
// if the value is 'off', it returns false
// if it is unset, it returns the default value
// otherwise, an error is returned.
func (e *envParamsImpl) GetOnOff(key string, setters ...GetEnvSetter) (on bool, err error) {
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
		return false, fmt.Errorf("environment variable %q value is %q and can only be \"on\" or \"off\"", key, s)
	}
}

// GetValueIfInside obtains the value stored for a named environment variable if it is part of a
// list of possible values. You can optionally specify a defaultValue
func (e *envParamsImpl) GetValueIfInside(key string, possibilities []string, setters ...GetEnvSetter) (value string, err error) {
	s, err := e.GetEnv(key, setters...)
	if err != nil {
		return "", err
	}
	for _, possibility := range possibilities {
		if s == possibility {
			return s, nil
		}
	}
	return "", fmt.Errorf("environment variable %q value is %q and can only be one of: %s", key, s, strings.Join(possibilities, ", "))
}

// GetDuration gets the duration from the environment variable corresponding to the key provided.
// If none is set, it returns defaultValue * timeUnit.
func (e *envParamsImpl) GetDuration(key string, timeUnit time.Duration, setters ...GetEnvSetter) (duration time.Duration, err error) {
	s, err := e.GetEnv(key, setters...)
	if err != nil {
		return 0, err
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("environment variable %q duration value %q is not a valid integer", key, s)
	} else if value == 0 {
		return 0, fmt.Errorf("environment variable %q duration value cannot be 0", key)
	} else if value < 0 {
		return 0, fmt.Errorf("environment variable %q duration value cannot be lower than 0", key)
	}
	return time.Duration(value) * timeUnit, nil
}

// GetHTTPTimeout returns the HTTP client timeout duration in milliseconds
// from the environment variable HTTP_TIMEOUT
func (e *envParamsImpl) GetHTTPTimeout(timeUnit time.Duration, setters ...GetEnvSetter) (timeout time.Duration, err error) {
	return e.GetDuration("HTTP_TIMEOUT", timeUnit, setters...)
}

// GetUserID obtains the user ID running the program
func (e *envParamsImpl) GetUserID() int {
	return e.getuid()
}

// GetGroupID obtains the user ID running the program
func (e *envParamsImpl) GetGroupID() int {
	return e.getgid()
}

// GetListeningPort obtains and checks the listening port
// from the environment variable LISTENING_PORT
func (e *envParamsImpl) GetListeningPort(setters ...GetEnvSetter) (listeningPort string, err error) {
	setters = append([]GetEnvSetter{Default("8000")}, setters...)
	listeningPort, err = e.GetEnv("LISTENING_PORT", setters...)
	if err != nil {
		return listeningPort, err
	}
	uid := e.getuid()
	warning, err := verifyListeningPort(listeningPort, uid)
	if warning != "" {
		logging.Warn(warning)
	}
	return listeningPort, err
}

// GetRootURL obtains and checks the root URL
// from the environment variable ROOT_URL
func (e *envParamsImpl) GetRootURL(setters ...GetEnvSetter) (rootURL string, err error) {
	setters = append([]GetEnvSetter{Default("/")}, setters...)
	rootURL, err = e.GetEnv("ROOT_URL", setters...)
	if err != nil {
		return rootURL, err
	}
	rootURL = path.Clean(rootURL)
	if err := verifyRootURL(rootURL); err != nil {
		return "", fmt.Errorf("environment variable ROOT_URL: %w", err)
	}
	rootURL = strings.TrimSuffix(rootURL, "/") // already have / from paths of router
	return rootURL, nil
}

// GetDatabaseDetails obtains the SQL database details from the
// environment variables SQL_USER, SQL_PASSWORD and SQL_DBNAME
func (e *envParamsImpl) GetDatabaseDetails() (hostname, user, password, dbName string, err error) {
	hostname, err = e.GetEnv("SQL_HOST", Default("postgres"))
	if err != nil {
		return hostname, user, password, dbName, err
	}
	if err := verifyHostname(hostname); err != nil {
		return hostname, user, password, dbName,
			fmt.Errorf("Postgres parameters: %w", err)
	}
	user, err = e.GetEnv("SQL_USER", Default("postgres"))
	if err != nil {
		return hostname, user, password, dbName, err
	}
	password, err = e.GetEnv("SQL_PASSWORD", Default("postgres"))
	if err != nil {
		return hostname, user, password, dbName, err
	}
	dbName, err = e.GetEnv("SQL_DBNAME", Default("postgres"))
	if err != nil {
		return hostname, user, password, dbName, err
	}
	// TODO port
	return hostname, user, password, dbName, nil
}

// GetRedisDetails obtains the Redis details from the
// environment variables REDIS_HOST, REDIS_PORT and REDIS_PASSWORD
func (e *envParamsImpl) GetRedisDetails() (hostname, port, password string, err error) {
	hostname, err = e.GetEnv("REDIS_HOST", Default("redis"))
	if err != nil {
		return hostname, port, password, err
	}
	if err := verifyHostname(hostname); err != nil {
		return hostname, port, password,
			fmt.Errorf("environment variable \"REDIS_HOST\": %w", err)
	}
	port, err = e.GetEnv("REDIS_PORT", Default("6379"))
	if err != nil {
		return hostname, port, password, err
	}
	if err := verification.VerifyPort(port); err != nil {
		return hostname, port, password,
			fmt.Errorf("environment variable REDIS_PORT: %w", err)
	}
	password, err = e.GetEnv("REDIS_PASSWORD")
	if err != nil {
		return hostname, port, password, err
	}
	return hostname, port, password, nil
}

// GetExeDir obtains the executable directory
func (e *envParamsImpl) GetExeDir() (dir string, err error) {
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
// It uses defaultValue if no value is found
func (e *envParamsImpl) GetPath(key string, setters ...GetEnvSetter) (path string, err error) {
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
// using the environment variables LOG_ENCODING, LOG_LEVEL and NODE_ID.
func (e *envParamsImpl) GetLoggerConfig() (encoding string, level logging.Level, nodeID int, err error) {
	encoding, err = e.GetLoggerEncoding()
	if err != nil {
		return "", 0, 0, fmt.Errorf("logger configuration error: %w", err)
	}
	level, err = e.GetLoggerLevel()
	if err != nil {
		return "", 0, 0, fmt.Errorf("logger configuration error: %w", err)
	}
	nodeID, err = e.GetNodeID()
	if err != nil {
		return "", 0, 0, fmt.Errorf("logger configuration error: %w", err)
	}
	return encoding, level, nodeID, nil
}

// GetLoggerEncoding obtains the logging encoding
// from the environment variable LOG_ENCODING
func (e *envParamsImpl) GetLoggerEncoding(setters ...GetEnvSetter) (encoding string, err error) {
	setters = append(setters, Default("json"))
	s, err := e.GetEnv("LOG_ENCODING", setters...)
	if err != nil {
		return "", err
	}
	s = strings.ToLower(s)
	if s != "json" && s != "console" {
		return "", fmt.Errorf("environment variable LOG_ENCODING: logger encoding %q unrecognized", s)
	}
	return s, nil
}

// GetLoggerLevel obtains the logging level
// from the environment variable LOG_LEVEL
func (e *envParamsImpl) GetLoggerLevel(setters ...GetEnvSetter) (level logging.Level, err error) {
	setters = append(setters, Default("json"))
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

// GetNodeID obtains the node instance ID from the environment variable
// NODE_ID
func (e *envParamsImpl) GetNodeID(setters ...GetEnvSetter) (nodeID int, err error) {
	setters = append(setters, Default("0"))
	s, err := e.GetEnv("NODE_ID", setters...)
	if err != nil {
		return nodeID, err
	}
	nodeID, err = strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("environment variable NODE_ID: %w", err)
	}
	return nodeID, nil
}

// GetURL obtains the URL for the environment variable for the key given.
// It returns the URL of defaultValue if defaultValue is not empty.
// If no defaultValue is given, it returns nil.
func (e *envParamsImpl) GetURL(key string, setters ...GetEnvSetter) (URL *url.URL, err error) {
	s, err := e.GetEnv(key, setters...)
	if s == "" {
		return nil, nil
	}
	URL, err = url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("environment variable %q URL value: %w", key, err)
	}
	return URL, nil
}

// GetGotifyURL obtains the URL for Gotify server
// from the environment variable GOTIFY_URL.
// It returns a nil URL if no value is found.
func (e *envParamsImpl) GetGotifyURL(setters ...GetEnvSetter) (URL *url.URL, err error) {
	return e.GetURL("GOTIFY_URL", setters...)
}

// GetGotifyToken obtains the token for the app on the Gotify server
// from the environment variable GOTIFY_TOKEN.
func (e *envParamsImpl) GetGotifyToken(setters ...GetEnvSetter) (token string, err error) {
	return e.GetEnv("GOTIFY_TOKEN", setters...)
}
