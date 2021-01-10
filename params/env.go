package params

import (
	"errors"
	"fmt"
	"net"
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
	Port(key string, optionSetters ...OptionSetter) (port uint16, err error)
	ListeningPort(key string, optionSetters ...OptionSetter) (port uint16, warning string, err error)
	ListeningAddress(key string, optionSetters ...OptionSetter) (address, warning string, err error)
	RootURL(key string, optionSetters ...OptionSetter) (rootURL string, err error)
	Path(key string, optionSetters ...OptionSetter) (path string, err error)
	LoggerEncoding(optionSetters ...OptionSetter) (encoding logging.Encoding, err error)
	LoggerLevel(optionSetters ...OptionSetter) (level logging.Level, err error)
	URL(key string, optionSetters ...OptionSetter) (URL *liburl.URL, err error)
}

type envParams struct {
	getuid func() int
	getenv func(key string) string
	regex  verification.Regex
	unset  func(k string) error
	fpAbs  func(s string) (string, error)
}

// NewEnv returns a new Env object.
func NewEnv() Env {
	return &envParams{
		getuid: os.Getuid,
		getenv: os.Getenv,
		regex:  verification.NewRegex(),
		unset:  os.Unsetenv,
		fpAbs:  filepath.Abs,
	}
}

type envOptions struct {
	compulsory         bool
	caseSensitiveValue bool
	unset              bool
	defaultValue       string
	retroKeys          []string
	onRetro            func(oldKey, newKey string)
}

// OptionSetter is a setter for options to Get functions.
type OptionSetter func(options *envOptions) error

// Compulsory forces the environment variable to contain a value.
func Compulsory() OptionSetter {
	return func(options *envOptions) error {
		if len(options.defaultValue) > 0 {
			return fmt.Errorf("cannot make environment variable value compulsory with a default value")
		}
		options.compulsory = true
		return nil
	}
}

// Default sets a default string value for the environment variable if no value is found.
func Default(defaultValue string) OptionSetter {
	return func(options *envOptions) error {
		if options.compulsory {
			return fmt.Errorf("cannot set default value for environment variable value which is compulsory")
		}
		options.defaultValue = defaultValue
		return nil
	}
}

// CaseSensitiveValue makes the value processing case sensitive.
func CaseSensitiveValue() OptionSetter {
	return func(options *envOptions) error {
		options.caseSensitiveValue = true
		return nil
	}
}

// Unset unsets the environment variable after it has been read.
func Unset() OptionSetter {
	return func(options *envOptions) error {
		options.unset = true
		return nil
	}
}

// RetroKeys tries to read from retroactive environment variable keys
// and runs the function onRetro if any retro environment variable is not
// empty. RetroKeys overrides previous RetroKeys optionSetters passed.
func RetroKeys(keys []string, onRetro func(oldKey, newKey string)) OptionSetter {
	return func(options *envOptions) error {
		options.retroKeys = keys
		options.onRetro = onRetro
		return nil
	}
}

// Get returns the value stored for a named environment variable,
// and a default if no value is found.
func (e *envParams) Get(key string, optionSetters ...OptionSetter) (value string, err error) {
	options := envOptions{}
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
	options := envOptions{}
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
	options := envOptions{}
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

var ErrInvalidPort = errors.New("invalid port")

func (e *envParams) ListeningAddress(key string, optionSetters ...OptionSetter) (
	address, warning string, err error) {
	hostport, err := e.Get(key, optionSetters...)
	if err != nil {
		return "", "", err
	}
	host, port, err := net.SplitHostPort(hostport)
	if err != nil {
		return "", "", err
	}
	p, err := strconv.Atoi(port)
	if err != nil {
		return "", "", fmt.Errorf("%w: %s", ErrInvalidPort, err)
	}
	warning, err = e.checkListeningPort(uint16(p))
	if err != nil {
		return "", warning, fmt.Errorf("%w: %s", ErrInvalidPort, err)
	}
	address = host + ":" + port
	return address, warning, nil
}

// ListeningPort obtains and checks a port from an environment variable
// and returns a warning depending on the port and user ID running the program.
func (e *envParams) ListeningPort(key string, optionSetters ...OptionSetter) (port uint16, warning string, err error) {
	port, err = e.Port(key, optionSetters...)
	if err != nil {
		return 0, "", err
	}
	warning, err = e.checkListeningPort(port)
	return port, warning, err
}

var ErrReservedListeningPort = errors.New(
	"listening port cannot be in the reserved system ports range (1 to 1023) when running without root")

func (e *envParams) checkListeningPort(port uint16) (warning string, err error) {
	const (
		maxPrivilegedPort = 1023
		minDynamicPort    = 49151
	)
	if port <= maxPrivilegedPort {
		switch e.getuid() {
		case 0:
			warning = "listening port " +
				strconv.Itoa(int(port)) +
				" allowed to be in the reserved system ports range as you are running as root"
		case -1:
			warning = "listening port " +
				strconv.Itoa(int(port)) +
				" allowed to be in the reserved system ports range as you are running in Windows"
		default:
			err = fmt.Errorf("%w: port %d", ErrReservedListeningPort, port)
		}
	} else if port > minDynamicPort {
		// dynamic and/or private ports.
		warning = "listening port " +
			strconv.Itoa(int(port)) +
			" is in the dynamic/private ports range (above 49151)"
	}
	return warning, err
}

// RootURL obtains and checks the root URL from the environment variable specified by key.
func (e *envParams) RootURL(key string, optionSetters ...OptionSetter) (rootURL string, err error) {
	optionSetters = append([]OptionSetter{Default("/")}, optionSetters...)
	rootURL, err = e.Get(key, optionSetters...)
	if err != nil {
		return rootURL, err
	}
	rootURL = path.Clean(rootURL)
	if !e.regex.MatchRootURL(rootURL) {
		return "", fmt.Errorf("environment variable ROOT_URL value %q is not valid", rootURL)
	}
	rootURL = strings.TrimSuffix(rootURL, "/") // already have / from paths of router
	return rootURL, nil
}

var ErrInvalidPath = errors.New("invalid filepath")

// Path obtains a path from the environment variable corresponding
// to key, and verifies it is valid. If it is a relative path,
// it is converted to an absolute path.
func (e *envParams) Path(key string, optionSetters ...OptionSetter) (path string, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return "", err
	}
	path, err = e.fpAbs(s)
	if err != nil {
		return "", fmt.Errorf(
			"%w: for environment variable %s: %s",
			ErrInvalidPath, key, err)
	}
	return path, nil
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
