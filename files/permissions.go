package files

import (
	"fmt"
	"os"

	"github.com/phayes/permbits"
)

// GetUserPermissions obtains the permissions of the user owning the file
func (f *fileManager) GetUserPermissions(filePath string) (read, write, execute bool, err error) {
	info, err := f.fileStat(filePath)
	if err != nil {
		return false, false, false, fmt.Errorf("cannot get permissions: %w", err)
	}
	permissions := permbits.FileMode(info.Mode())
	return permissions.UserRead(), permissions.UserWrite(), permissions.UserExecute(), nil
}

func (f *fileManager) SetUserPermissions(filepath string, mod os.FileMode) error {
	exists, err := f.FileExists(filepath)
	if err != nil {
		return fmt.Errorf("cannot set user permissions: %w", err)
	} else if !exists {
		return fmt.Errorf("file %q does not exist", filepath)
	}
	return f.chmod(filepath, mod)
}
