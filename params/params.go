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

// GetEnv returns the value stored for a named environment variable,
// and a default if no value is found
func GetEnv(key, defaultValue string) (value string) {
	value = os.Getenv(key)
	if len(value) == 0 {
		value = defaultValue
	}
	return value
}

// GetEnvInt returns the value stored for a named environment variable,
// and a default if no value is found. If the value is not a valid
// integer, an error is returned.
func GetEnvInt(key string, defaultValue int) (n int, err error) {
	s := GetEnv(key, "")
	if s == "" {
		return defaultValue, nil
	}
	n, err = strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("environment variable %q: %w", key, err)
	}
	return n, nil
}

// GetEnvIntRange returns the value stored for a named environment variable,
// and a default if no value is found. If the value is not a valid
// integer in the range specified, an error is returned.
func GetEnvIntRange(key string, lower, upper, defaultValue int) (n int, err error) {
	n = defaultValue
	s := GetEnv(key, "")
	if s != "" {
		n, err = strconv.Atoi(s)
		if err != nil {
			return 0, fmt.Errorf("environment variable %q: %w", key, err)
		}
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
func GetYesNo(key string, defaultValue bool) (yes bool, err error) {
	s := GetEnv(key, "")
	switch s {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	case "":
		return defaultValue, nil
	default:
		return false, fmt.Errorf("environment variable %q value is %q and can only be \"yes\" or \"no\"", key, s)
	}
}

// GetOnOff obtains the value stored for a named environment variable and returns:
// if the value is 'on', it returns true
// if the value is 'off', it returns false
// if it is unset, it returns the default value
// otherwise, an error is returned.
func GetOnOff(key string, defaultValue bool) (on bool, err error) {
	s := GetEnv(key, "")
	switch s {
	case "on":
		return true, nil
	case "off":
		return false, nil
	case "":
		return defaultValue, nil
	default:
		return false, fmt.Errorf("environment variable %q value is %q and can only be \"on\" or \"off\"", key, s)
	}
}

// GetValueIfInside obtains the value stored for a named environment variable if it is part of a
// list of possible values. You can optionally specify a defaultValue
func GetValueIfInside(key string, possibilities []string, compulsory bool, defaultValue string) (value string, err error) {
	s := GetEnv(key, "")
	if !compulsory && s == "" {
		s = defaultValue
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
func GetDuration(key string, defaultValue int, timeUnit time.Duration) (duration time.Duration, err error) {
	s := GetEnv(key, "")
	if s == "" {
		return time.Duration(defaultValue) * timeUnit, nil
	}
	value, err := strconv.Atoi(s)
	if err != nil {
		return duration, fmt.Errorf("environment variable %q duration: %w", key, err)
	} else if value == 0 {
		return duration, fmt.Errorf("environment variable %q duration cannot be 0", key)
	} else if value < 0 {
		return duration, fmt.Errorf("environment variable %q duration cannot be lower than 0", key)
	}
	return time.Duration(value) * timeUnit, nil
}

// GetHTTPTimeout returns the HTTP client timeout duration in milliseconds
// from the environment variable HTTP_TIMEOUT
func GetHTTPTimeout(defaultMilliseconds int) (duration time.Duration, err error) {
	duration, err = GetDuration("HTTP_TIMEOUT", defaultMilliseconds, time.Millisecond)
	if err != nil {
		return duration, err
	}
	return duration, nil
}

// GetUserID obtains the user ID running the program
func GetUserID() int {
	return os.Geteuid()
}

// GetGroupID obtains the user ID running the program
func GetGroupID() int {
	return os.Getegid()
}

// GetListeningPort obtains and checks the listening port
// from the environment variable LISTENING_PORT
func GetListeningPort() (listeningPort string, err error) {
	listeningPort = GetEnv("LISTENING_PORT", "8000")
	uid := GetUserID()
	warning, err := verifyListeningPort(listeningPort, uid)
	if warning != "" {
		logging.Warn(warning)
	}
	return listeningPort, err
}

// GetRootURL obtains and checks the root URL
// from the environment variable ROOT_URL
func GetRootURL() (rootURL string, err error) {
	rootURL = GetEnv("ROOT_URL", "/")
	rootURL = path.Clean(rootURL)
	if err := verifyRootURL(rootURL); err != nil {
		return "", fmt.Errorf("environment variable ROOT_URL: %w", err)
	}
	rootURL = strings.TrimSuffix(rootURL, "/") // already have / from paths of router
	return rootURL, nil
}

// GetDatabaseDetails obtains the SQL database details from the
// environment variables SQL_USER, SQL_PASSWORD and SQL_DBNAME
func GetDatabaseDetails() (hostname, user, password, dbName string, err error) {
	hostname = GetEnv("SQL_HOST", "postgres")
	if err := verifyHostname(hostname); err != nil {
		return hostname, user, password, dbName,
			fmt.Errorf("Postgres parameters: %w", err)
	}
	// TODO port
	return hostname,
		GetEnv("SQL_USER", "postgres"),
		GetEnv("SQL_PASSWORD", "postgres"),
		GetEnv("SQL_DBNAME", "postgres"),
		nil
}

// GetRedisDetails obtains the Redis details from the
// environment variables REDIS_HOST, REDIS_PORT and REDIS_PASSWORD
func GetRedisDetails() (hostname, port, password string, err error) {
	hostname = GetEnv("REDIS_HOST", "redis")
	if err := verifyHostname(hostname); err != nil {
		return hostname, port, password,
			fmt.Errorf("environment variable REDIS_HOST: %w", err)
	}
	port = GetEnv("REDIS_PORT", "6379")
	if err := verification.VerifyPort(port); err != nil {
		return hostname, port, password,
			fmt.Errorf("environment variable REDIS_PORT: %w", err)
	}
	return hostname, port,
		GetEnv("REDIS_PASSWORD", ""),
		nil
}

// GetExeDir obtains the executable directory
func GetExeDir() (dir string, err error) {
	ex, err := os.Executable()
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
func GetPath(key, defaultValue string) (path string, err error) {
	s := GetEnv(key, defaultValue)
	if !filepath.IsAbs(s) {
		exDir, err := GetExeDir()
		if err != nil {
			return "", err
		}
		s = filepath.Join(exDir, s)
	}
	return filepath.Abs(s)
}

// GetLoggerConfig obtains configuration details for the logger
// using the environment variables LOG_ENCODING, LOG_LEVEL and NODE_ID.
func GetLoggerConfig() (encoding string, level logging.Level, nodeID int, err error) {
	encoding, err = GetLoggerEncoding()
	if err != nil {
		return "", 0, 0, fmt.Errorf("logger configuration error: %w", err)
	}
	level, err = GetLoggerLevel()
	if err != nil {
		return "", 0, 0, fmt.Errorf("logger configuration error: %w", err)
	}
	nodeID, err = GetNodeID()
	if err != nil {
		return "", 0, 0, fmt.Errorf("logger configuration error: %w", err)
	}
	return encoding, level, nodeID, nil
}

// GetLoggerEncoding obtains the logging encoding
// from the environment variable LOG_ENCODING
func GetLoggerEncoding() (encoding string, err error) {
	s := GetEnv("LOG_ENCODING", "json")
	s = strings.ToLower(s)
	if s != "json" && s != "console" {
		return "", fmt.Errorf("environment variable LOG_ENCODING: logger encoding %q unrecognized", s)
	}
	return s, nil
}

// GetLoggerLevel obtains the logging level
// from the environment variable LOG_LEVEL
func GetLoggerLevel() (level logging.Level, err error) {
	s := GetEnv("LOG_LEVEL", "info")
	switch strings.ToLower(s) {
	case "info":
		return logging.InfoLevel, nil
	case "warning":
		return logging.WarnLevel, nil
	case "error":
		return logging.ErrorLevel, nil
	case "":
		return logging.InfoLevel, nil
	default:
		return level, fmt.Errorf("environment variable LOG_LEVEL: logger level %q unrecognized", s)
	}
}

// GetNodeID obtains the node instance ID from the environment variable
// NODE_ID
func GetNodeID() (nodeID int, err error) {
	s := GetEnv("NODE_ID", "0")
	nodeID, err = strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("environment variable NODE_ID: %w", err)
	}
	return nodeID, nil
}

// GetURL obtains the URL for the environment variable for the key given.
// It returns the URL of defaultValue if defaultValue is not empty.
// If no defaultValue is given, it returns nil.
func GetURL(key, defaultValue string) (URL *url.URL, err error) {
	s := GetEnv(key, defaultValue)
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
func GetGotifyURL() (*url.URL, error) {
	return GetURL("GOTIFY_URL", "")
}

// GetGotifyToken obtains the token for the app on the Gotify server
// from the environment variable GOTIFY_TOKEN.
func GetGotifyToken() (token string, err error) {
	token = GetEnv("GOTIFY_TOKEN", "")
	if token == "" {
		return "", fmt.Errorf("environment variable GOTIFY_TOKEN value not provided")
	}
	return token, nil
}
