package params

import (
	"errors"
	"fmt"
	"net"
	"strconv"
)

var ErrInvalidPort = errors.New("invalid port")

func (e *Env) ListeningAddress(key string, optionSetters ...OptionSetter) (
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
