package params

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"
)

var ErrRootURLNotValid = errors.New("root URL is not valid")

var rootURLRegex = regexp.MustCompile(`^\/[a-zA-Z0-9\-_/\+]*$`)

// RootURL obtains and checks the root URL from the environment variable specified by key.
func (e *Env) RootURL(key string, optionSetters ...OptionSetter) (rootURL string, err error) {
	optionSetters = append([]OptionSetter{Default("/")}, optionSetters...)
	rootURL, err = e.Get(key, optionSetters...)
	if err != nil {
		return rootURL, err
	}
	rootURL = path.Clean(rootURL)
	if !rootURLRegex.MatchString(rootURL) {
		return "", fmt.Errorf("%w: %s", ErrRootURLNotValid, rootURL)
	}
	rootURL = strings.TrimSuffix(rootURL, "/") // already have / from paths of router
	return rootURL, nil
}
