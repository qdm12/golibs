package files

import (
	"io/fs"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_FileManager_GetUserPermissions(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		perms   fs.FileMode
		read    bool
		write   bool
		execute bool
	}{
		"no_permission": {},
		"read_permission": {
			perms: 0o0400,
			read:  true,
		},
		"write_permission": {
			perms: 0o0200,
			write: true,
		},
		"execute_permission": {
			perms:   0o0100,
			execute: true,
		},
		"read_write_permissions": {
			perms: 0o0600,
			read:  true,
			write: true,
		},
		"read_write_execute_permissions": {
			perms:   0o0700,
			read:    true,
			write:   true,
			execute: true,
		},
		"group_and_other_ignored": {
			perms: 0o0024,
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

			read, write, execute, err := manager.GetUserPermissions(file.Name())

			require.NoError(t, err)
			assert.Equal(t, testCase.read, read)
			assert.Equal(t, testCase.write, write)
			assert.Equal(t, testCase.execute, execute)
		})
	}
}

func Test_FileManager_GetGroupPermissions(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		perms   fs.FileMode
		read    bool
		write   bool
		execute bool
	}{
		"no_permission": {},
		"read_permission": {
			perms: 0o0040,
			read:  true,
		},
		"write_permission": {
			perms: 0o0020,
			write: true,
		},
		"execute_permission": {
			perms:   0o0010,
			execute: true,
		},
		"read_write_permissions": {
			perms: 0o0060,
			read:  true,
			write: true,
		},
		"read_write_execute_permissions": {
			perms:   0o0070,
			read:    true,
			write:   true,
			execute: true,
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

			read, write, execute, err := manager.GetGroupPermissions(file.Name())

			require.NoError(t, err)
			assert.Equal(t, testCase.read, read)
			assert.Equal(t, testCase.write, write)
			assert.Equal(t, testCase.execute, execute)
		})
	}
}

func Test_FileManager_GetOthersPermissions(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		perms   fs.FileMode
		read    bool
		write   bool
		execute bool
	}{
		"no_permission": {},
		"read_permission": {
			perms: 0o0004,
			read:  true,
		},
		"write_permission": {
			perms: 0o0002,
			write: true,
		},
		"execute_permission": {
			perms:   0o0001,
			execute: true,
		},
		"read_write_permissions": {
			perms: 0o0006,
			read:  true,
			write: true,
		},
		"read_write_execute_permissions": {
			perms:   0o0007,
			read:    true,
			write:   true,
			execute: true,
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

			read, write, execute, err := manager.GetOthersPermissions(file.Name())

			require.NoError(t, err)
			assert.Equal(t, testCase.read, read)
			assert.Equal(t, testCase.write, write)
			assert.Equal(t, testCase.execute, execute)
		})
	}
}
