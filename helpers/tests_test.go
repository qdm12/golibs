package helpers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AssertErrosEqual(t *testing.T) {
	t.Parallel()
	tests := map[string]struct {
		expected error
		actual   error
		success  bool
	}{
		"empty errors": {
			errors.New(""),
			errors.New(""),
			true,
		},
		"nil errors": {
			nil,
			nil,
			true,
		},
		"nil error not equal to empty error": {
			nil,
			errors.New(""),
			false,
		},
		"empty error not equal to nil error": {
			errors.New(""),
			nil,
			false,
		},
		"errors with different content": {
			errors.New("abc"),
			errors.New("abx"),
			false,
		},
	}
	for name, tc := range tests {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			tMock := &testing.T{}
			out := AssertErrorsEqual(tMock, tc.expected, tc.actual)
			assert.Equal(t, tc.success, out)
		})
	}
}
