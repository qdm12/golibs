package params

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrNotAnInteger = errors.New("value is not an integer")
)

// Int returns the value stored for a named environment variable,
// and a default if no value is found. If the value is not a valid
// integer, an error is returned.
func (e *Env) Int(key string, optionSetters ...OptionSetter) (n int, err error) {
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
func (e *Env) IntRange(key string, lower, upper int, optionSetters ...OptionSetter) (n int, err error) {
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
