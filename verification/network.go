package verification

import (
	"fmt"
	"strconv"
)

// VerifyPort verifies a port number string is valid
func VerifyPort(port string) error {
	value, err := strconv.Atoi(port)
	if err != nil {
		return fmt.Errorf("port %q is not a valid integer", port)
	} else if value < 1 {
		return fmt.Errorf("port %s cannot be lower than 1", port)
	} else if value > 65535 {
		return fmt.Errorf("port %s cannot be higher than 65535", port)
	}
	return nil
}
