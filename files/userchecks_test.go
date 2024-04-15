package files

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FileManager_isAccessible(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		perms         fs.FileMode
		uid           int
		gid           int
		accessibility accessible
		ok            bool
		errMessage    string
	}{
		// "no_permission": {
		// 	accessibility: readable,
		// },
		"readable_other_read": {
			accessibility: readable,
			perms:         0o0004,
			ok:            true,
		},
		"readable_group_read": {
			accessibility: readable,
			perms:         0o0040,
			ok:            true,
		},
		"readable_user_read": {
			accessibility: readable,
			perms:         0o0400,
			ok:            true,
		},
		"readable_user_write": {
			accessibility: readable,
			perms:         0o0200,
		},
		"writable_user_read": {
			accessibility: writable,
			perms:         0o0400,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			file, err := os.CreateTemp(t.TempDir(), "")
			require.NoError(t, err)
			err = os.Chmod(file.Name(), testCase.perms)
			require.NoError(t, err)

			manager := NewFileManager()

			ok, err := manager.isAccessible(file.Name(),
				testCase.uid, testCase.gid, testCase.accessibility)

			if testCase.errMessage != "" {
				require.Error(t, err)
				assert.EqualError(t, err, testCase.errMessage)
			} else {
				require.NoError(t, err)
			}
			assert.Equal(t, testCase.ok, ok)
		})
	}
}
