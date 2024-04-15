package files

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"strings"
)

// WriteLinesToFile writes a slice of strings as lines to a file.
// It creates any directory not existing in the file path if necessary.
func (f *FileManager) WriteLinesToFile(filePath string, lines []string, options ...WriteOptionSetter) error {
	var data []byte
	if len(lines) > 0 && (len(lines) != 1 || len(lines[0]) > 0) {
		data = []byte(strings.Join(lines, "\n"))
	}
	return f.WriteToFile(filePath, data, options...)
}

// Touch creates an empty file at the file path given.
// It creates any directory not existing in the file path if necessary.
func (f *FileManager) Touch(filePath string, options ...WriteOptionSetter) error {
	return f.WriteToFile(filePath, nil, options...)
}

func (f *FileManager) CreateDir(filePath string, setters ...WriteOptionSetter) error {
	const defaultPermissions os.FileMode = 0700
	options := newWriteOptions(defaultPermissions)
	for _, setter := range setters {
		setter(&options)
	}

	exists, err := f.FilepathExists(filePath)
	if err != nil {
		return err
	}
	if exists { //nolint:nestif
		isFile, err := f.IsFile(filePath)
		if err != nil {
			return err
		} else if isFile {
			return fmt.Errorf("%w: %s", fs.ErrExist, filePath)
		}
		err = f.SetUserPermissions(filePath, options.permissions)
		if err != nil {
			return fmt.Errorf("setting file permissions: %w", err)
		}
	} else {
		err = f.mkdirAll(filePath, options.permissions)
		if err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
	}
	if options.ownership != nil {
		err = f.chown(filePath, options.ownership.UID, options.ownership.GID)
		if err != nil {
			return fmt.Errorf("setting file permissions: %w", err)
		}
	}
	return nil
}

var ErrIsDirectory = errors.New("path is a directory")

// WriteToFile writes data bytes to a file, and creates any
// directory not existing in the file path if necessary.
func (f *FileManager) WriteToFile(filePath string, data []byte,
	setters ...WriteOptionSetter) (err error) {
	const defaultPermissions os.FileMode = 0600
	options := newWriteOptions(defaultPermissions)
	for _, setter := range setters {
		setter(&options)
	}
	exists, err := f.FileExists(filePath)
	if err != nil {
		return fmt.Errorf("checking if file exists: %w", err)
	}

	if exists {
		isDir, err := f.IsDirectory(filePath)
		if err != nil {
			return fmt.Errorf("checking if path is a directory: %w", err)
		} else if isDir {
			return fmt.Errorf("%w: %s", ErrIsDirectory, filePath)
		}
	} else {
		parentDir := f.filepathDir(filePath)
		err = f.mkdirAll(parentDir, 0700)
		if err != nil {
			return fmt.Errorf("creating parent directory: %w", err)
		}
	}

	err = f.writeFile(filePath, data, options.permissions)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}

	if options.ownership != nil {
		err = f.chown(filePath, options.ownership.UID, options.ownership.GID)
		if err != nil {
			return fmt.Errorf("changing ownership of %s to %d:%d: %w",
				filePath, options.ownership.UID, options.ownership.GID, err)
		}
	}

	return nil
}
