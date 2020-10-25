package verification

import (
	"fmt"
	"strconv"
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
		return 0, fmt.Errorf("port %q is not a valid integer", s)
	case portInt < minPort:
		return 0, fmt.Errorf("port %s cannot be lower than 1", s)
	case portInt > maxPort:
		return 0, fmt.Errorf("port %s cannot be higher than 65535", s)
	default:
		return uint16(portInt), nil
	}
}
