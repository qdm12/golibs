package params

import (
	"errors"
	"fmt"
	"strings"
)

var ErrNotOneOf = errors.New("value is not within the accepted values")

// Inside obtains the value stored for a named environment variable if it is part of a
// list of possible values. You can optionally specify a defaultValue.
func (e *Env) Inside(key string, possibilities []string, optionSetters ...OptionSetter) (
	value string, err error) {
	options := envOptions{}
	for _, setter := range optionSetters {
		_ = setter(&options) // error is checked in e.Get
	}
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return "", err
	} else if len(s) == 0 && !options.compulsory {
		return "", nil
	}
	for _, possibility := range possibilities {
		if options.caseSensitiveValue && s == possibility {
			return s, nil
		} else if !options.caseSensitiveValue && strings.EqualFold(s, possibility) {
			return strings.ToLower(s), nil
		}
	}
	csvPossibilities := strings.Join(possibilities, ", ")
	return "", fmt.Errorf("%w: %s: it can only be one of: %s", ErrNotOneOf, s, csvPossibilities)
}
