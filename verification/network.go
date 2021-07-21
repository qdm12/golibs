package verification

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrPortIsNotInteger = errors.New("port is not an integer")
	ErrPortLowerThanOne = errors.New("port cannot be lower than 1")
	ErrPortTooHigh      = errors.New("port cannot be higher than 65535")
)

func (v *verifier) VerifyPort(port string) error {
	_, err := ParsePort(port)
	return err
}

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
