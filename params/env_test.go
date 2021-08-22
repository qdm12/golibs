package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewEnv(t *testing.T) {
	t.Parallel()
	e := NewEnv()
	assert.NotNil(t, e)
}
