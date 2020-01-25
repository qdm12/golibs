package files

import (
	"fmt"

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
