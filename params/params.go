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

	"github.com/qdm12/golibs/verification"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	return strconv.Atoi(s)
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
		return duration, fmt.Errorf("duration %q: %w", key, err)
	} else if value == 0 {
		return duration, fmt.Errorf("duration %q cannot be 0", key)
	} else if value < 0 {
		return duration, fmt.Errorf("duration %q cannot be lower than 0", key)
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

// GetListeningPort obtains and checks the listening port
// from the environment variable LISTENING_PORT
func GetListeningPort() (listeningPort string, err error) {
	listeningPort = GetEnv("LISTENING_PORT", "8000")
	uid := os.Geteuid()
	warning, err := verifyListeningPort(listeningPort, uid)
	zap.L().Warn(warning)
	return listeningPort, err
}

// GetRootURL obtains and checks the root URL
// from the environment variable ROOT_URL
func GetRootURL() (rootURL string, err error) {
	rootURL = GetEnv("ROOT_URL", "/")
	rootURL = path.Clean(rootURL)
	if err := verifyRootURL(rootURL); err != nil {
		return "", err
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
			fmt.Errorf("Redis parameters: %w", err)
	}
	port = GetEnv("REDIS_PORT", "6379")
	if err := verification.VerifyPort(port); err != nil {
		return hostname, port, password,
			fmt.Errorf("Redis parameters: %w", err)
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
// to key, and verifies it is valid. It uses defaultValue if no value is
// found
func GetPath(key, defaultValue string) (path string, err error) {
	s := GetEnv(key, defaultValue)
	return filepath.Abs(s)
}

// GetLoggerEncoding obtains the logging encoding for Zap
// from the environment variable LOG_ENCODING
func GetLoggerEncoding() (encoding string, err error) {
	s := GetEnv("LOG_ENCODING", "json")
	s = strings.ToLower(s)
	if s != "json" && s != "console" {
		return "", fmt.Errorf("logger encoding %q unrecognized", s)
	}
	return s, nil
}

// GetLoggerLevel obtains the logging level for Zap
// from the environment variable LOG_LEVEL
func GetLoggerLevel() (level zapcore.Level, err error) {
	s := GetEnv("LOG_LEVEL", "info")
	switch strings.ToLower(s) {
	case "info":
		return zap.InfoLevel, nil
	case "warning":
		return zap.WarnLevel, nil
	case "error":
		return zap.ErrorLevel, nil
	case "":
		return zap.InfoLevel, nil
	default:
		return level, fmt.Errorf("logger level %q unrecognized", s)
	}
}

// GetNodeID obtains the node instance ID from the environment variable
// NODE_ID
func GetNodeID() (nodeID int, err error) {
	s := GetEnv("NODE_ID", "0")
	nodeID, err = strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("Node ID %q is not a valid integer", nodeID)
	}
	return nodeID, nil
}

// GetURL obtains the URL for the environment variable for the key given.
// It returns the URL of defaultValue if defaultValue is not empty.
// If no defaultValue is given, it returns nil.
func GetURL(key, defaultValue string) (*url.URL, error) {
	URL := GetEnv(key, defaultValue)
	if URL == "" {
		return nil, nil
	}
	return url.Parse(URL)
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
		return "", fmt.Errorf("Gotify token not provided")
	}
	return token, nil
}
