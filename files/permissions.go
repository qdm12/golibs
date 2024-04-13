package files

import (
	"fmt"
	"os"

	"github.com/phayes/permbits"
)

// GetUserPermissions obtains the permissions of the user owning the file.
func (f *FileManager) GetUserPermissions(filePath string) (read, write, execute bool, err error) {
	info, err := f.fileStat(filePath)
	if err != nil {
		return false, false, false, fmt.Errorf("cannot get permissions: %w", err)
	}
	permissions := permbits.FileMode(info.Mode())
	return permissions.UserRead(), permissions.UserWrite(), permissions.UserExecute(), nil
}

// GetGroupPermissions obtains the permissions of the group owning the file path.
func (f *FileManager) GetGroupPermissions(filePath string) (read, write, execute bool, err error) {
	info, err := f.fileStat(filePath)
	if err != nil {
		return false, false, false, fmt.Errorf("cannot get permissions: %w", err)
	}
	permissions := permbits.FileMode(info.Mode())
	return permissions.GroupRead(), permissions.GroupWrite(), permissions.GroupExecute(), nil
}

// GetOthersPermissions obtains the permissions for users and groups not owning the file path.
func (f *FileManager) GetOthersPermissions(filePath string) (read, write, execute bool, err error) {
	info, err := f.fileStat(filePath)
	if err != nil {
		return false, false, false, fmt.Errorf("cannot get permissions: %w", err)
	}
	permissions := permbits.FileMode(info.Mode())
	return permissions.OtherRead(), permissions.OtherWrite(), permissions.OtherExecute(), nil
}

func (f *FileManager) SetUserPermissions(filepath string, mod os.FileMode) error {
	exists, err := f.FileExists(filepath)
	if err != nil {
		return err
	} else if !exists {
		return ErrFileNotExist
	}
	return f.chmod(filepath, mod)
}
