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
	CSV(key string, optionSetters ...OptionSetter) (values []string, err error)
	CSVInside(key string, possibilities []string, optionSetters ...OptionSetter) (values []string, err error)
	Duration(key string, optionSetters ...OptionSetter) (duration time.Duration, err error)
	Port(key string, optionSetters ...OptionSetter) (port uint16, err error)
	ListeningPort(key string, optionSetters ...OptionSetter) (port uint16, warning string, err error)
	ListeningAddress(key string, optionSetters ...OptionSetter) (address, warning string, err error)
	RootURL(key string, optionSetters ...OptionSetter) (rootURL string, err error)
	Path(key string, optionSetters ...OptionSetter) (path string, err error)
	LogCaller(key string, optionSetters ...OptionSetter) (caller logging.Caller, err error)
	LogLevel(key string, optionSetters ...OptionSetter) (level logging.Level, err error)
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

var (
	ErrCompulsoryConflictWithDefault = errors.New("cannot make environment variable value compulsory with a default value")
	ErrDefaultConflictWithCompulsory = errors.New("cannot set a default for a compulsory environment variable value")
)

// Compulsory forces the environment variable to contain a value.
func Compulsory() OptionSetter {
	return func(options *envOptions) error {
		if len(options.defaultValue) > 0 {
			return ErrCompulsoryConflictWithDefault
		}
		options.compulsory = true
		return nil
	}
}

// Default sets a default string value for the environment variable if no value is found.
func Default(defaultValue string) OptionSetter {
	return func(options *envOptions) error {
		if options.compulsory {
			return ErrDefaultConflictWithCompulsory
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

var (
	ErrOption  = errors.New("option error")
	ErrNoValue = errors.New("no value found")
)

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
			return "", fmt.Errorf("%w: %s", ErrOption, err)
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
			return "", ErrNoValue
		}
		value = options.defaultValue
	}
	if !options.caseSensitiveValue {
		value = strings.ToLower(value)
	}
	return value, nil
}

var (
	ErrNotAnInteger = errors.New("value is not an integer")
)

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
		return 0, fmt.Errorf("%w: %s", ErrNotAnInteger, s)
	}
	return n, nil
}

var (
	ErrNotInRange = errors.New("value is not in range")
)

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
		return 0, fmt.Errorf("%w: %s", ErrNotAnInteger, s)
	}
	if n < lower || n > upper {
		return 0, fmt.Errorf("%w: %d is not between %d and %d", ErrNotInRange, n, lower, upper)
	}
	return n, nil
}

var ErrNotYesNo = errors.New("value can only be 'yes' or 'no'")

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
		return false, fmt.Errorf("%w: %s", ErrNotYesNo, s)
	}
}

var ErrNotOnOff = errors.New("value can only be 'on' or 'off'")

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
		return false, fmt.Errorf("%w: %s", ErrNotOnOff, s)
	}
}

var ErrNotOneOf = errors.New("value is not within the accepted values")

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
	return "", fmt.Errorf("%w: %s: it can only be one of: %s", ErrNotOneOf, s, csvPossibilities)
}

func (e *envParams) CSV(key string, optionSetters ...OptionSetter) (values []string, err error) {
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
	return strings.Split(csv, ","), nil
}

var ErrInvalidValueFound = errors.New("at least one value is not within the accepted values")

func (e *envParams) CSVInside(key string, possibilities []string, optionSetters ...OptionSetter) (
	values []string, err error) {
	values, err = e.CSV(key, optionSetters...)
	if err != nil {
		return nil, err
	} else if values == nil {
		return nil, nil
	}

	options := envOptions{}
	for _, setter := range optionSetters {
		_ = setter(&options) // error is checked in e.Get
	}
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
		csvInvalidMessages := strings.Join(invalidMessages, ", ")
		csvPossibilities := strings.Join(possibilities, ", ")
		return nil, fmt.Errorf("%w: invalid values found: %s; possible values are: %s",
			ErrInvalidValueFound, csvInvalidMessages, csvPossibilities)
	}
	return values, nil
}

var (
	ErrDurationMalformed = errors.New("duration is malformed")
	ErrDurationNegative  = errors.New("duration is negative")
)

// Duration gets the duration from the environment variable corresponding to the key provided.
func (e *envParams) Duration(key string, optionSetters ...OptionSetter) (duration time.Duration, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return 0, err
	} else if len(s) == 0 {
		return 0, nil
	}
	duration, err = time.ParseDuration(s)
	switch {
	case err != nil:
		return 0, fmt.Errorf("%w: %s: %s", ErrDurationMalformed, s, err)
	case duration < 0:
		return 0, fmt.Errorf("%w: %s", ErrDurationNegative, duration.String())
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

var ErrRootURLNotValid = errors.New("root URL is not valid")

// RootURL obtains and checks the root URL from the environment variable specified by key.
func (e *envParams) RootURL(key string, optionSetters ...OptionSetter) (rootURL string, err error) {
	optionSetters = append([]OptionSetter{Default("/")}, optionSetters...)
	rootURL, err = e.Get(key, optionSetters...)
	if err != nil {
		return rootURL, err
	}
	rootURL = path.Clean(rootURL)
	if !e.regex.MatchRootURL(rootURL) {
		return "", fmt.Errorf("%w: %s", ErrRootURLNotValid, rootURL)
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
		return "", fmt.Errorf("%w: %s: %s", ErrInvalidPath, path, err)
	}
	return path, nil
}

var ErrUnknownLogCaller = errors.New("unknown log caller")

func (e *envParams) LogCaller(key string, optionSetters ...OptionSetter) (caller logging.Caller, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return caller, err
	}
	switch strings.ToLower(s) {
	case "hidden":
		return logging.CallerHidden, nil
	case "short":
		return logging.CallerShort, nil
	}
	return caller, fmt.Errorf("%w: %s: can be one of: hidden, short", ErrUnknownLogCaller, s)
}

var ErrUnknownLogLevel = errors.New("unknown log level")

// LogLevel obtains the log level from an environment variable.
func (e *envParams) LogLevel(key string, optionSetters ...OptionSetter) (level logging.Level, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return level, err
	}
	switch strings.ToLower(s) {
	case "debug":
		return logging.LevelDebug, nil
	case "info":
		return logging.LevelInfo, nil
	case "warning":
		return logging.LevelWarn, nil
	case "error":
		return logging.LevelError, nil
	default:
		return level, fmt.Errorf("%w: %s: can be one of: debug, info, warning, error",
			ErrUnknownLogLevel, s)
	}
}

var (
	ErrURLNotValid = errors.New("url is not valid")
	ErrURLNotHTTP  = errors.New("url is not http(s)")
)

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
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %s", ErrURLNotValid, s, err)
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, fmt.Errorf("%w: %s", ErrURLNotHTTP, url.String())
	}
	return url, nil
}
