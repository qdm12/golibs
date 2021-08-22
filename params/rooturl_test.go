package params

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RootURL(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		rootURL       string
		err           error
	}{
		"key with valid value": {
			envValue: "/a",
			rootURL:  "/a",
		},
		"key with valid value and trail /": {
			envValue: "/a/",
			rootURL:  "/a",
		},
		"key with invalid value": {
			envValue: "/a?",
			err:      fmt.Errorf("root URL is not valid: /a?"),
		},
		"key without value": {},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("/a")},
			rootURL:       "/a",
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           fmt.Errorf("option error: cannot make environment variable value compulsory with a default value"),
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &Env{
				getenv: func(key string) string {
					return tc.envValue
				},
			}
			rootURL, err := e.RootURL("ROOT_URL", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.rootURL, rootURL)
		})
	}
}
