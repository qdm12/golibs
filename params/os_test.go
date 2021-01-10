package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UID(t *testing.T) {
	t.Parallel()
	const expectedUID = 1
	o := &osImpl{
		getuid: func() int {
			return expectedUID
		},
	}
	uid := o.UID()
	assert.Equal(t, expectedUID, uid)
}

func Test_GID(t *testing.T) {
	t.Parallel()
	const expectedUID = 1
	o := &osImpl{
		getgid: func() int {
			return expectedUID
		},
	}
	gid := o.GID()
	assert.Equal(t, expectedUID, gid)
}
