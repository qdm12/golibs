package params

import (
	"errors"
	"fmt"
	"strings"
)

func (e *Env) CSV(key string, optionSetters ...OptionSetter) (values []string, err error) {
	options := envOptions{}
	for _, setter := range optionSetters {
		_ = setter(&options) // error is checked in e.Get
	}
	csv, err := e.Get(key, optionSetters...)
	if err != nil {
		return nil, err
	}
	if !options.compulsory && len(csv) == 0 {
		return nil, nil
	}
	return strings.Split(csv, ","), nil
}

var ErrInvalidValueFound = errors.New("at least one value is not within the accepted values")

func (e *Env) CSVInside(key string, possibilities []string, optionSetters ...OptionSetter) (
	values []string, err error) {
	values, err = e.CSV(key, optionSetters...)
	if err != nil {
		return nil, err
	} else if values == nil {
		return nil, nil
	}

	options := envOptions{}
	for _, setter := range optionSetters {
		_ = setter(&options) // error is checked in e.Get
	}
	type valuePosition struct {
		position int
		value    string
	}
	var invalidValues []valuePosition
	for i, value := range values {
		found := false
		for _, possibility := range possibilities {
			if options.caseSensitiveValue {
				if value == possibility {
					found = true
					break
				}
			} else {
				if strings.EqualFold(value, possibility) {
					values[i] = strings.ToLower(value)
					found = true
					break
				}
			}
		}
		if !found {
			invalidValues = append(invalidValues, valuePosition{i + 1, value})
		}
	}
	if L := len(invalidValues); L > 0 {
		invalidMessages := make([]string, L)
		for i := range invalidValues {
			invalidMessages[i] = fmt.Sprintf("value %q at position %d", invalidValues[i].value, invalidValues[i].position)
		}
		csvInvalidMessages := strings.Join(invalidMessages, ", ")
		csvPossibilities := strings.Join(possibilities, ", ")
		return nil, fmt.Errorf("%w: invalid values found: %s; possible values are: %s",
			ErrInvalidValueFound, csvInvalidMessages, csvPossibilities)
	}
	return values, nil
}
