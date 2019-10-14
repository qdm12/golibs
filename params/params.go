package params

import (
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"fmt"

	"github.com/qdm12/golibs/verification"
)

func getEnv(key, defaultValue string) (value string) {
	value = os.Getenv(key)
	if len(value) == 0 {
		value = defaultValue
	}
	return value
}

// GetListeningPort obtains and checks the listening port from Viper (env variable or config file, etc.)
func GetListeningPort() (listeningPort, warning string, err error) {
	listeningPort = getEnv("LISTENINGPORT", "8000")
	uid := os.Geteuid()
	warning, err = verifyListeningPort(listeningPort, uid)
	return listeningPort, warning, err
}

// GetRootURL obtains and checks the root URL from Viper (env variable or config file, etc.)
func GetRootURL() (rootURL string, err error) {
	rootURL = getEnv("ROOTURL", "/")
	if err := verifyRootURL(rootURL); err != nil {
		return rootURL, err
	}
	rootURL = strings.ReplaceAll(rootURL, "//", "/")
	rootURL = strings.TrimSuffix(rootURL, "/") // already have / from paths of router
	return rootURL, nil
}

// GetDatabaseDetails obtains the SQL database details from Viper (env variable or config file, etc.)
func GetDatabaseDetails() (hostname, user, password, dbName string, err error) {
	hostname = getEnv("sqlhost", "postgres")
	if err := verifyHostname(hostname); err != nil {
		return hostname, user, password, dbName,
			fmt.Errorf("Postgres parameters: %w", err)
	}
	// TODO port
	return hostname,
		getEnv("sqluser", "postgres"),
		getEnv("sqlpassword", "postgres"),
		getEnv("sqldbname", "postgres"),
		nil
}

// GetRedisDetails obtains the Redis details from Viper (env variable or config file, etc.)
func GetRedisDetails() (hostname, port, password string, err error) {
	hostname = getEnv("redishost", "redis")
	if err := verifyHostname(hostname); err != nil {
		return hostname, port, password,
			fmt.Errorf("Redis parameters: %w", err)
	}
	port = getEnv("redisport", "6379")
	if err := verification.VerifyPort(port); err != nil {
		return hostname, port, password,
			fmt.Errorf("Redis parameters: %w", err)
	}
	return hostname, port,
		getEnv("redispassword", ""),
		nil
}

// GetDir obtains the executable directory
func GetDir() (dir string, err error) {
	ex, err := os.Executable()
	if err != nil {
		return dir, err
	}
	dir = filepath.Dir(ex)
	return dir, nil
}

// GetLoggerMode obtains the logging mode from Viper (env variable or config file, etc.)
func GetLoggerMode() (production bool, err error) {
	s := getEnv("logging", "production")
	switch strings.ToLower(s) {
	case "production":
		return true, nil
	case "development":
		return false, nil
	}
	return false, fmt.Errorf("logging mode %s unrecognized", s)
}

// GetGotifyURL obtains the URL to the Gotify server
func GetGotifyURL() (*url.URL, error) {
	URL := getEnv("gotifyurl", "")
	if URL == "" {
		return nil, nil
	}
	return url.Parse(URL)
}

// GetGotifyToken obtains the token for the app on the Gotify server
func GetGotifyToken() (token string, err error) {
	token = getEnv("gotifytoken", "")
	if token == "" {
		return "", fmt.Errorf("Gotify token not provided")
	}
	return token, nil
}
