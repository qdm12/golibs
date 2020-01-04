package helpers

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_AssertErrorsEqual(t *testing.T) {
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
			tMock := &mockTestingT{}
			if !tc.success {
				tMock.On("Errorf", mock.AnythingOfType("string"), mock.AnythingOfType("string"))
			}
			if tc.expected != nil && tc.actual == nil {
				tMock.On("FailNow")
			}
			out := AssertErrorsEqual(tMock, tc.expected, tc.actual)
			assert.Equal(t, tc.success, out)
		})
	}
}
