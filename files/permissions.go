package files

import (
	"fmt"
	"os"
)

// GetUserPermissions obtains the permissions of the user owning the file.
func (f *FileManager) GetUserPermissions(filePath string) (read, write, execute bool, err error) {
	info, err := f.fileStat(filePath)
	if err != nil {
		return false, false, false, fmt.Errorf("cannot get permissions: %w", err)
	}
	mode := info.Mode()
	perm := mode.Perm()
	return perm&0400 != 0, perm&0200 != 0, perm&0100 != 0, nil
}

// GetGroupPermissions obtains the permissions of the group owning the file path.
func (f *FileManager) GetGroupPermissions(filePath string) (read, write, execute bool, err error) {
	info, err := f.fileStat(filePath)
	if err != nil {
		return false, false, false, fmt.Errorf("cannot get permissions: %w", err)
	}
	mode := info.Mode()
	perm := mode.Perm()
	return perm&0040 != 0, perm&0020 != 0, perm&0010 != 0, nil
}

// GetOthersPermissions obtains the permissions for users and groups not owning the file path.
func (f *FileManager) GetOthersPermissions(filePath string) (read, write, execute bool, err error) {
	info, err := f.fileStat(filePath)
	if err != nil {
		return false, false, false, fmt.Errorf("cannot get permissions: %w", err)
	}
	mode := info.Mode()
	perm := mode.Perm()
	return perm&0004 != 0, perm&0002 != 0, perm&0001 != 0, nil
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
