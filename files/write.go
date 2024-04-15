package files

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// WriteLinesToFile writes a slice of strings as lines to a file.
// It creates any directory not existing in the file path if necessary.
func WriteLinesToFile(filePath string, lines []string, options ...WriteOptionSetter) error {
	var data []byte
	if len(lines) > 0 && (len(lines) != 1 || len(lines[0]) > 0) {
		data = []byte(strings.Join(lines, "\n"))
	}
	return WriteToFile(filePath, data, options...)
}

// Touch creates an empty file at the file path given.
// It creates any directory not existing in the file path if necessary.
func Touch(filePath string, options ...WriteOptionSetter) error {
	return WriteToFile(filePath, nil, options...)
}

func CreateDir(filePath string, setters ...WriteOptionSetter) error {
	const defaultPermissions os.FileMode = 0700
	options := newWriteOptions(defaultPermissions)
	for _, setter := range setters {
		setter(&options)
	}

	exists, err := FilepathExists(filePath)
	if err != nil {
		return err
	}
	if exists { //nolint:nestif
		isFile, err := IsFile(filePath)
		if err != nil {
			return err
		} else if isFile {
			return fmt.Errorf("%w: %s", fs.ErrExist, filePath)
		}
		err = os.Chmod(filePath, options.permissions)
		if err != nil {
			return fmt.Errorf("setting file permissions: %w", err)
		}
	} else {
		err = os.MkdirAll(filePath, options.permissions)
		if err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
	}
	if options.ownership != nil {
		err = os.Chown(filePath, options.ownership.UID, options.ownership.GID)
		if err != nil {
			return fmt.Errorf("setting file permissions: %w", err)
		}
	}
	return nil
}

var ErrIsDirectory = errors.New("path is a directory")

// WriteToFile writes data bytes to a file, and creates any
// directory not existing in the file path if necessary.
func WriteToFile(filePath string, data []byte,
	setters ...WriteOptionSetter) (err error) {
	const defaultPermissions os.FileMode = 0600
	options := newWriteOptions(defaultPermissions)
	for _, setter := range setters {
		setter(&options)
	}
	exists, err := FileExists(filePath)
	if err != nil {
		return fmt.Errorf("checking if file exists: %w", err)
	}

	if exists {
		isDir, err := IsDirectory(filePath)
		if err != nil {
			return fmt.Errorf("checking if path is a directory: %w", err)
		} else if isDir {
			return fmt.Errorf("%w: %s", ErrIsDirectory, filePath)
		}
	} else {
		parentDir := filepath.Dir(filePath)
		err = CreateDir(parentDir)
		if err != nil {
			return fmt.Errorf("creating parent directory: %w", err)
		}
	}

	err = os.WriteFile(filePath, data, options.permissions)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	if options.ownership != nil {
		err = os.Chown(filePath, options.ownership.UID, options.ownership.GID)
		if err != nil {
			return fmt.Errorf("changing ownership of %s to %d:%d: %w",
				filePath, options.ownership.UID, options.ownership.GID, err)
		}
	}

	return nil
}
