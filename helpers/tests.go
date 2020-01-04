package helpers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Helper()
}

// AssertErrorsEqual asserts errors match, so that they are both nil or their messages
// are both equal
func AssertErrorsEqual(t testingT, expected, actual error) (success bool) {
	t.Helper()
	if expected != nil {
		require.Error(t, actual)
		if actual == nil {
			return false
		}
		return assert.Equal(t, expected.Error(), actual.Error())
	}
	return assert.NoError(t, actual)
}
