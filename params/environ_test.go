package params

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parseEnviron(t *testing.T) {
	t.Parallel()

	t.Log(os.Environ())

	environ := []string{
		"",
		"keyA",
		"keyB=value",
		"keyC=value=value2",
		"keyD=",
	}

	kv := parseEnviron(environ)

	expectedKV := map[string]string{
		"":     "",
		"keyA": "",
		"keyB": "value",
		"keyC": "value=value2",
		"keyD": "",
	}

	assert.Equal(t, expectedKV, kv)
}
