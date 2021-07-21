package redis

import (
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var (
	ErrSendCommand   = errors.New("cannot send command to output buffer")
	ErrFlushCommands = errors.New("cannot flush commands")
	ErrReceiveResult = errors.New("cannot receive result")
	ErrDoCommand     = errors.New("cannot Do command")
	ErrResultNotOk   = errors.New("result is not OK")
)

// IsNil returns true if the reply is nil.
func IsNil(reply interface{}) bool {
	switch reply.(type) { //nolint:gocritic
	case nil:
		return true
	}
	return false
}

// CheckOKString returns an error if the reply is not the "OK" string.
func CheckOKString(reply interface{}) error {
	ok, err := redis.String(reply, nil)
	if err != nil {
		return err
	} else if ok != "OK" {
		return fmt.Errorf("%w: %s", ErrResultNotOk, ok)
	}
	return nil
}
