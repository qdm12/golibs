package params

import (
	"errors"
	"fmt"
	"strconv"
)

// Port obtains and checks a port number from the
// environment variable corresponding to the key provided.
func (e *Env) Port(key string, optionSetters ...OptionSetter) (port uint16, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return 0, err
	}
	return ParsePort(s)
}

var (
	ErrPortIsNotInteger = errors.New("port is not an integer")
	ErrPortLowerThanOne = errors.New("port cannot be lower than 1")
	ErrPortTooHigh      = errors.New("port cannot be higher than 65535")
)

// ParsePort verifies a port number string is valid and
// returns the port as an uint16.
func ParsePort(s string) (port uint16, err error) {
	const minPort = 1
	const maxPort = 65535
	portInt, err := strconv.Atoi(s)
	switch {
	case err != nil:
		return 0, fmt.Errorf("%w: %s", ErrPortIsNotInteger, s)
	case portInt < minPort:
		return 0, fmt.Errorf("%w: %d", ErrPortLowerThanOne, portInt)
	case portInt > maxPort:
		return 0, fmt.Errorf("%w: %d", ErrPortTooHigh, portInt)
	default:
		return uint16(portInt), nil
	}
}
