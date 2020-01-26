package format

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ArgsToString(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		args []interface{}
		s    string
	}{
		"no args":        {nil, ""},
		"one arg":        {[]interface{}{2}, "2"},
		"two args":       {[]interface{}{2, 3}, "2 3"},
		"format and arg": {[]interface{}{"one is %d", 1}, "one is 1"},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := ArgsToString(tc.args...)
			assert.Equal(t, tc.s, s)
		})
	}
}

func Test_ReadableBytes(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		bytes []byte
		s     string
	}{
		"nil":            {nil, ""},
		"empty":          {[]byte{}, ""},
		"one byte":       {[]byte{1}, "2"},
		"multiple bytes": {[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, "4HUtbHhN2TkpR"},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s := ReadableBytes(tc.bytes)
			assert.Equal(t, tc.s, s)
		})
	}
}
