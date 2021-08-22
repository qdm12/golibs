package params

import (
	"errors"
	"fmt"
	"strings"

	"github.com/qdm12/golibs/logging"
)

var ErrUnknownLogCaller = errors.New("unknown log caller")

func (e *Env) LogCaller(key string, optionSetters ...OptionSetter) (caller logging.Caller, err error) {
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
func (e *Env) LogLevel(key string, optionSetters ...OptionSetter) (level logging.Level, err error) {
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
