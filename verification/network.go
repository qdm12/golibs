package verification

import (
	"fmt"
	"strconv"
)

// VerifyPort verifies a port number string is valid
func VerifyPort(port string) error {
	const minPort = 1
	const maxPort = 65535
	value, err := strconv.Atoi(port)
	switch {
	case err != nil:
		return fmt.Errorf("port %q is not a valid integer", port)
	case value < minPort:
		return fmt.Errorf("port %s cannot be lower than 1", port)
	case value > maxPort:
		return fmt.Errorf("port %s cannot be higher than 65535", port)
	default:
		return nil
	}
}
