package params

import (
	"errors"
	"fmt"
	"strconv"
)

// ListeningPort obtains and checks a port from an environment variable
// and returns a warning depending on the port and user ID running the program.
func (e *Env) ListeningPort(key string, optionSetters ...OptionSetter) (port uint16, warning string, err error) {
	port, err = e.Port(key, optionSetters...)
	if err != nil {
		return 0, "", err
	}
	warning, err = e.checkListeningPort(port)
	return port, warning, err
}

var ErrReservedListeningPort = errors.New(
	"listening port cannot be in the reserved system ports range (1 to 1023) when running without root")

func (e *Env) checkListeningPort(port uint16) (warning string, err error) {
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
