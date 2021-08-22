package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_New(t *testing.T) {
	t.Parallel()
	env := New()
	assert.NotNil(t, env)
	assert.NotNil(t, env.kv)
}
