package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/phayes/permbits"
)

// FilepathExists returns true if a fiel path exists.
func FilepathExists(filePath string) (exists bool, err error) {
	_, err = os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FileExists returns true if a file exists.
func FileExists(filePath string) (exists bool, err error) {
	exists, err = FilepathExists(filePath)
	if err != nil {
		return false, err
	} else if !exists {
		return false, nil
	}
	info, _ := os.Stat(filePath)
	if info.IsDir() {
		return false, nil
	}
	return true, nil
}

// DirectoryExists returns true if a directory exists.
func DirectoryExists(filePath string) (exists bool, err error) {
	exists, err = FilepathExists(filePath)
	if err != nil {
		return false, err
	} else if !exists {
		return false, nil
	}
	info, _ := os.Stat(filePath)
	if !info.IsDir() {
		return false, nil
	}
	return true, nil
}

// GetOwnership obtains the user ID and group ID owning the file or directory
// or returns 0 and 0 if running in Windows (no IDs)
func GetOwnership(filePath string) (userID, groupID int, err error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot get ownership: %w", err)
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0, nil // Windows
	}
	return int(stat.Uid), int(stat.Gid), nil
}

// GetUserPermissions obtains the permissions of the user owning the file
func GetUserPermissions(filePath string) (read, write, execute bool, err error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, false, false, fmt.Errorf("cannot get permissions: %w", err)
	}
	permissions := permbits.FileMode(info.Mode())
	return permissions.UserRead(), permissions.UserWrite(), permissions.UserExecute(), nil
}

// WriteLinesToFile writes a slice of strings as lines to a file.
// It creates any directory not existing in the file path if necessary.
func WriteLinesToFile(filePath string, lines []string) error {
	data := []byte(strings.Join(lines, "\n"))
	return WriteToFile(filePath, data)
}

// Touch creates an empty file at the file path given.
// It creates any directory not existing in the file path if necessary.
func Touch(filePath string) error {
	return WriteToFile(filePath, nil)
}

// WriteToFile writes data bytes to a file, and creates any
// directory not existing in the file path if necessary.
func WriteToFile(filePath string, data []byte) error {
	exists, err := FileExists(filePath)
	if err != nil {
		return fmt.Errorf("cannot write to file: %w", err)
	} else if !exists {
		parentDir := filepath.Dir(filePath)
		err := os.MkdirAll(parentDir, 0700)
		if err != nil {
			return fmt.Errorf("cannot write to file %q: %w", filePath, err)
		}
	}
	return ioutil.WriteFile(filePath, data, 0600)
}

// ReadFile reads the entire data of a file.
func ReadFile(filePath string) (data []byte, err error) {
	exists, err := FileExists(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	} else if !exists {
		return nil, fmt.Errorf("cannot read file: %q does not exist", filePath)
	}
	return ioutil.ReadFile(filePath)
}
