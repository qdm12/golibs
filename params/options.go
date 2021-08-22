package params

import "errors"

type envOptions struct {
	compulsory         bool
	caseSensitiveValue bool
	unset              bool
	defaultValue       string
	retroKeys          []string
	onRetro            func(oldKey, newKey string)
}

// OptionSetter is a setter for options to Get functions.
type OptionSetter func(options *envOptions) error

var (
	ErrCompulsoryConflictWithDefault = errors.New("cannot make environment variable value compulsory with a default value")
	ErrDefaultConflictWithCompulsory = errors.New("cannot set a default for a compulsory environment variable value")
)

// Compulsory forces the environment variable to contain a value.
func Compulsory() OptionSetter {
	return func(options *envOptions) error {
		if len(options.defaultValue) > 0 {
			return ErrCompulsoryConflictWithDefault
		}
		options.compulsory = true
		return nil
	}
}

// Default sets a default string value for the environment variable if no value is found.
func Default(defaultValue string) OptionSetter {
	return func(options *envOptions) error {
		if options.compulsory {
			return ErrDefaultConflictWithCompulsory
		}
		options.defaultValue = defaultValue
		return nil
	}
}

// CaseSensitiveValue makes the value processing case sensitive.
func CaseSensitiveValue() OptionSetter {
	return func(options *envOptions) error {
		options.caseSensitiveValue = true
		return nil
	}
}

// Unset unsets the environment variable after it has been read.
func Unset() OptionSetter {
	return func(options *envOptions) error {
		options.unset = true
		return nil
	}
}

// RetroKeys tries to read from retroactive environment variable keys
// and runs the function onRetro if any retro environment variable is not
// empty. RetroKeys overrides previous RetroKeys optionSetters passed.
func RetroKeys(keys []string, onRetro func(oldKey, newKey string)) OptionSetter {
	return func(options *envOptions) error {
		options.retroKeys = keys
		options.onRetro = onRetro
		return nil
	}
}
