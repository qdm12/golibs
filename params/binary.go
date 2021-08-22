package params

import (
	"errors"
	"fmt"
)

var ErrNotYesNo = errors.New("value can only be 'yes' or 'no'")

// YesNo obtains the value stored for a named environment variable and returns:
// if the value is 'yes', it returns true
// if the value is 'no', it returns false
// if it is unset, it returns the default value
// otherwise, an error is returned.
func (e *Env) YesNo(key string, optionSetters ...OptionSetter) (yes bool, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return false, err
	}
	switch s {
	case "yes":
		return true, nil
	case "no":
		return false, nil
	default:
		return false, fmt.Errorf("%w: %s", ErrNotYesNo, s)
	}
}

var ErrNotOnOff = errors.New("value can only be 'on' or 'off'")

// OnOff obtains the value stored for a named environment variable and returns:
// if the value is 'on', it returns true
// if the value is 'off', it returns false
// if it is unset, it returns the default value
// otherwise, an error is returned.
func (e *Env) OnOff(key string, optionSetters ...OptionSetter) (on bool, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return false, err
	}
	switch s {
	case "on":
		return true, nil
	case "off":
		return false, nil
	default:
		return false, fmt.Errorf("%w: %s", ErrNotOnOff, s)
	}
}
