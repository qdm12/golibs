package params

import (
	"errors"
	"fmt"
)

var ErrInvalidPath = errors.New("invalid filepath")

// Path obtains a path from the environment variable corresponding
// to key, and verifies it is valid. If it is a relative path,
// it is converted to an absolute path.
func (e *Env) Path(key string, optionSetters ...OptionSetter) (path string, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return "", err
	}
	path, err = e.fpAbs(s)
	if err != nil {
		return "", fmt.Errorf("%w: %s: %s", ErrInvalidPath, path, err)
	}
	return path, nil
}
