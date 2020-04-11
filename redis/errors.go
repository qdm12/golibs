package redis

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

// WrapSendErr wraps an error with a context message of a send command
func WrapSendErr(err error) error {
	return fmt.Errorf("cannot send command to output buffer for Redis: %w", err)
}

// WrapSendErr wraps an error with a context message of a flush command
func WrapFlushErr(err error) error {
	return fmt.Errorf("cannot flush commands to Redis: %w", err)
}

// WrapSendErr wraps an error with a context message of a receive
func WrapReceiveErr(err error) error {
	return fmt.Errorf("cannot receive result from Redis: %w", err)
}

// WrapSendErr wraps an error with a context message of a do command
func WrapDoErr(err error) error {
	return fmt.Errorf("cannot Do command with Redis: %w", err)
}

// IsNil returns true if the reply is nil
func IsNil(reply interface{}) bool {
	switch reply.(type) { //nolint:gocritic
	case nil:
		return true
	}
	return false
}

// CheckString returns the string and an error if the reply is not a string
func CheckString(reply interface{}) (string, error) {
	s, err := redis.String(reply, nil)
	if err != nil {
		return "", fmt.Errorf("result from Redis is not a string: %w", err)
	}
	return s, nil
}

// CheckInteger returns the integer and an error if the reply is not an integer
func CheckInteger(reply interface{}) (int, error) {
	n, err := redis.Int(reply, nil)
	if err != nil {
		return 0, fmt.Errorf("result from Redis is not an integer: %w", err)
	}
	return n, nil
}

// CheckOKString returns an error if the reply is not the "OK" string
func CheckOKString(reply interface{}) error {
	ok, err := redis.String(reply, nil)
	if err != nil {
		return fmt.Errorf("result from Redis is not a string: %w", err)
	}
	if ok != "OK" {
		return fmt.Errorf("result from Redis is not OK")
	}
	return nil
}

// CheckSlice returns the slice of slice of bytes and an error if the reply is not slices of bytes
func CheckSlice(reply interface{}) ([][]byte, error) {
	byteSlices, err := redis.ByteSlices(reply, nil)
	if err != nil {
		return nil, fmt.Errorf("result from Redis is not byte slices: %w", err)
	}
	return byteSlices, nil
}

// CheckValues returns a slice of values and an error if the reply is not an array of results
func CheckValues(reply interface{}) ([]interface{}, error) {
	replies, err := redis.Values(reply, nil)
	if err != nil {
		return nil, fmt.Errorf("result from Redis is not an array of results: %w", err)
	}
	return replies, nil
}
