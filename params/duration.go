package params

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrDurationMalformed = errors.New("duration is malformed")
	ErrDurationNegative  = errors.New("duration is negative")
)

// Duration gets the duration from the environment variable corresponding to the key provided.
func (e *Env) Duration(key string, optionSetters ...OptionSetter) (duration time.Duration, err error) {
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
