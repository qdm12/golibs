package files

import (
	"fmt"
	"io/fs"
	"os"
)

// GetUserPermissions obtains the permissions of the user owning the file.
func (f *FileManager) GetUserPermissions(filePath string) (read, write, execute bool, err error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, false, false, err
	}
	mode := info.Mode()
	perm := mode.Perm()
	return perm&0400 != 0, perm&0200 != 0, perm&0100 != 0, nil
}

// GetGroupPermissions obtains the permissions of the group owning the file path.
func (f *FileManager) GetGroupPermissions(filePath string) (read, write, execute bool, err error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, false, false, err
	}
	mode := info.Mode()
	perm := mode.Perm()
	return perm&0040 != 0, perm&0020 != 0, perm&0010 != 0, nil
}

// GetOthersPermissions obtains the permissions for users and groups not owning the file path.
func (f *FileManager) GetOthersPermissions(filePath string) (read, write, execute bool, err error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, false, false, err
	}
	mode := info.Mode()
	perm := mode.Perm()
	return perm&0004 != 0, perm&0002 != 0, perm&0001 != 0, nil
}

func (f *FileManager) SetUserPermissions(filepath string, mod os.FileMode) error {
	exists, err := f.FileExists(filepath)
	if err != nil {
		return fmt.Errorf("checking file existence: %w", err)
	} else if !exists {
		return fmt.Errorf("%w: %s", fs.ErrNotExist, filepath)
	}
	return os.Chmod(filepath, mod)
}
