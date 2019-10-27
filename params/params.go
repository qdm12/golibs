package params

import (
	"net/url"
	"os"
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
// If none is set, it returns defaultValue.
func GetDuration(key string, defaultValue time.Duration) (duration time.Duration, err error) {
	value, err := GetEnvInt(key, -1)
	if err != nil {
		return duration, fmt.Errorf("duration %s: %w", key, err)
	} else if value == 0 {
		return duration, fmt.Errorf("duration %s cannot be 0", key)
	} else if value < 0 {
		return defaultValue, nil
	}
	return time.Duration(value) * time.Millisecond, nil
}

// GetHTTPTimeout returns the HTTP client timeout duration in milliseconds
// from the environment variable HTTPTIMEOUT
func GetHTTPTimeout(defaultValue time.Duration) (duration time.Duration, err error) {
	duration, err = GetDuration("HTTPTIMEOUT", defaultValue)
	if err != nil {
		return duration, err
	}
	return duration * time.Millisecond, nil
}

// GetListeningPort obtains and checks the listening port
// from the environment variable LISTENINGPORT
func GetListeningPort() (listeningPort string, err error) {
	listeningPort = GetEnv("LISTENINGPORT", "8000")
	uid := os.Geteuid()
	warning, err := verifyListeningPort(listeningPort, uid)
	zap.L().Warn(warning)
	return listeningPort, err
}

// GetRootURL obtains and checks the root URL
// from the environment variable ROOTURL
func GetRootURL() (rootURL string, err error) {
	rootURL = GetEnv("ROOTURL", "/")
	if err := verifyRootURL(rootURL); err != nil {
		return rootURL, err
	}
	rootURL = strings.ReplaceAll(rootURL, "//", "/")
	rootURL = strings.TrimSuffix(rootURL, "/") // already have / from paths of router
	return rootURL, nil
}

// GetDatabaseDetails obtains the SQL database details from the
// environment variables SQLUSER, SQLPASSWORD and SQLDBNAME
func GetDatabaseDetails() (hostname, user, password, dbName string, err error) {
	hostname = GetEnv("sqlhost", "postgres")
	if err := verifyHostname(hostname); err != nil {
		return hostname, user, password, dbName,
			fmt.Errorf("Postgres parameters: %w", err)
	}
	// TODO port
	return hostname,
		GetEnv("SQLUSER", "postgres"),
		GetEnv("SQLPASSWORD", "postgres"),
		GetEnv("SQLDBNAME", "postgres"),
		nil
}

// GetRedisDetails obtains the Redis details from the
// environment variables REDISHOST, REDISPORT and REDISPASSWORD
func GetRedisDetails() (hostname, port, password string, err error) {
	hostname = GetEnv("REDISHOST", "redis")
	if err := verifyHostname(hostname); err != nil {
		return hostname, port, password,
			fmt.Errorf("Redis parameters: %w", err)
	}
	port = GetEnv("redisport", "6379")
	if err := verification.VerifyPort(port); err != nil {
		return hostname, port, password,
			fmt.Errorf("Redis parameters: %w", err)
	}
	return hostname, port,
		GetEnv("redispassword", ""),
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
// from the environment variable LOGENCODING
func GetLoggerEncoding() (encoding string, err error) {
	s := GetEnv("LOGENCODING", "json")
	s = strings.ToLower(s)
	if s != "json" && s != "console" {
		return "", fmt.Errorf("logger encoding \"%s\" unrecognized", s)
	}
	return s, nil
}

// GetLoggerLevel obtains the logging level for Zap
// from the environment variable LOGLEVEL
func GetLoggerLevel() (level zapcore.Level, err error) {
	s := GetEnv("LOGLEVEL", "info")
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
		return level, fmt.Errorf("logger level \"%s\" unrecognized", s)
	}
}

// GetNodeID obtains the node instance ID from the environment variable
// NODEID
func GetNodeID() (nodeID int, err error) {
	s := GetEnv("NODEID", "0")
	nodeID, err = strconv.Atoi(s)
	if err != nil {
		return 0, fmt.Errorf("Node ID \"%d\" is not a valid integer", nodeID)
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
// from the environment variable GOTIFYURL.
// It returns a nil URL if no value is found.
func GetGotifyURL() (*url.URL, error) {
	return GetURL("GOTIFYURL", "")
}

// GetGotifyToken obtains the token for the app on the Gotify server
// from the environment variable GOTIFYTOKEN.
func GetGotifyToken() (token string, err error) {
	token = GetEnv("GOTIFYTOKEN", "")
	if token == "" {
		return "", fmt.Errorf("Gotify token not provided")
	}
	return token, nil
}
