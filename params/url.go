package params

import (
	"errors"
	"fmt"
	liburl "net/url"
)

var (
	ErrURLNotValid = errors.New("url is not valid")
	ErrURLNotHTTP  = errors.New("url is not http(s)")
)

// URL obtains the HTTP URL for the environment variable for the key given.
// It returns the URL of defaultValue if defaultValue is not empty.
// If no defaultValue is given, it returns nil.
func (e *Env) URL(key string, optionSetters ...OptionSetter) (url *liburl.URL, err error) {
	s, err := e.Get(key, optionSetters...)
	if err != nil {
		return nil, err
	} else if s == "" {
		return nil, nil
	}

	url, err = liburl.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("%w: %s: %s", ErrURLNotValid, s, err)
	}

	if url.Scheme != "http" && url.Scheme != "https" {
		return nil, fmt.Errorf("%w: %s", ErrURLNotHTTP, url.String())
	}
	return url, nil
}
