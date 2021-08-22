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

func Test_NewFromEnviron(t *testing.T) {
	t.Parallel()
	environ := []string{"key=value"}
	env := NewFromEnviron(environ)
	assert.Equal(t, map[string]string{"key": "value"}, env.kv)
	assert.NotNil(t, env)
}
