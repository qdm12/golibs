package params

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrOption  = errors.New("option error")
	ErrNoValue = errors.New("no value found")
)

// Get returns the value stored for a named environment variable,
// and a default if no value is found.
func (e *Env) Get(key string, optionSetters ...OptionSetter) (value string, err error) {
	options := envOptions{}
	defer func() {
		if options.unset {
			_ = e.unset(key)
			for _, retroKey := range options.retroKeys {
				_ = e.unset(retroKey)
			}
		}
	}()
	for _, setter := range optionSetters {
		if err := setter(&options); err != nil {
			return "", fmt.Errorf("%w: %s", ErrOption, err)
		}
	}
	value = e.getenv(key)
	if len(value) == 0 {
		for _, retroKey := range options.retroKeys {
			value = e.getenv(retroKey)
			if len(value) > 0 {
				options.onRetro(retroKey, key)
				break
			}
		}
	}
	if len(value) == 0 {
		if options.compulsory {
			return "", ErrNoValue
		}
		value = options.defaultValue
	}
	if !options.caseSensitiveValue {
		value = strings.ToLower(value)
	}
	return value, nil
}
