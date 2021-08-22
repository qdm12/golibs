package params

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_URL(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		envValue      string
		optionSetters []OptionSetter
		URL           string
		err           error
	}{
		"key with URL value": {"https://google.com", nil, "https://google.com", nil},
		"key with invalid value": {
			envValue: string([]byte{0}),
			err:      fmt.Errorf("url is not valid: \x00: parse \"\\x00\": net/url: invalid control character in URL")},
		"key with non HTTP value": {
			envValue: "google.com",
			err:      fmt.Errorf("url is not http(s): google.com"),
		},
		"key without value": {},
		"key without value and default": {
			optionSetters: []OptionSetter{Default("https://google.com")},
			URL:           "https://google.com",
		},
		"key without value and compulsory": {
			optionSetters: []OptionSetter{Compulsory()},
			err:           ErrNoValue,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			e := &Env{
				kv: map[string]string{"any": tc.envValue},
			}
			URL, err := e.URL("any", tc.optionSetters...)
			if tc.err != nil {
				require.Error(t, err)
				assert.Equal(t, tc.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}
			if URL == nil {
				assert.Empty(t, tc.URL)
			} else {
				assert.Equal(t, tc.URL, URL.String())
			}
		})
	}
}
