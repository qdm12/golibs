package files

import (
	"errors"
	"fmt"
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

var ErrFileExistsAtPath = errors.New("a file exists at this path")

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
	if exists {
		isFile, err := f.IsFile(filePath)
		if err != nil {
			return err
		} else if isFile {
			return ErrFileExistsAtPath
		}
		if err := f.SetUserPermissions(filePath, options.permissions); err != nil {
			return err
		}
	} else if err := f.mkdirAll(filePath, options.permissions); err != nil {
		return err
	}
	if options.ownership != nil {
		if err := f.chown(filePath, options.ownership.UID, options.ownership.GID); err != nil {
			return err
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
		return err
	}

	if options.ownership != nil {
		err = f.chown(filePath, options.ownership.UID, options.ownership.GID)
		if err != nil {
			return err
		}
	}

	return nil
}
