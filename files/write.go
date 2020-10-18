package files

import (
	"fmt"
	"os"
	"strings"
)

// WriteLinesToFile writes a slice of strings as lines to a file.
// It creates any directory not existing in the file path if necessary.
func (f *fileManager) WriteLinesToFile(filePath string, lines []string, options ...WriteOptionSetter) error {
	var data []byte
	if len(lines) > 0 && (len(lines) != 1 || len(lines[0]) > 0) {
		data = []byte(strings.Join(lines, "\n"))
	}
	return f.WriteToFile(filePath, data, options...)
}

// Touch creates an empty file at the file path given.
// It creates any directory not existing in the file path if necessary.
func (f *fileManager) Touch(filePath string, options ...WriteOptionSetter) error {
	return f.WriteToFile(filePath, nil, options...)
}

func (f *fileManager) CreateDir(filePath string, setters ...WriteOptionSetter) error {
	const defaultPermissions os.FileMode = 0700
	options := newWriteOptions(defaultPermissions)
	for _, setter := range setters {
		setter(&options)
	}
	errPrefix := fmt.Sprintf("cannot create directory %q: ", filePath)
	exists, err := f.FilepathExists(filePath)
	if err != nil {
		return fmt.Errorf("%s%w", errPrefix, err)
	}
	if exists {
		isFile, err := f.IsFile(filePath)
		if err != nil {
			return fmt.Errorf("%s%w", errPrefix, err)
		} else if isFile {
			return fmt.Errorf("%sa file exists at this path", errPrefix)
		}
		if err := f.SetUserPermissions(filePath, options.permissions); err != nil {
			return fmt.Errorf("%s%w", errPrefix, err)
		}
	} else if err := f.mkdirAll(filePath, options.permissions); err != nil {
		return err
	}
	if options.ownership != nil {
		if err := f.chown(filePath, options.ownership.UID, options.ownership.GID); err != nil {
			return fmt.Errorf("%s%w", errPrefix, err)
		}
	}
	return nil
}

// WriteToFile writes data bytes to a file, and creates any
// directory not existing in the file path if necessary.
func (f *fileManager) WriteToFile(filePath string, data []byte, setters ...WriteOptionSetter) error {
	const defaultPermissions os.FileMode = 0600
	options := newWriteOptions(defaultPermissions)
	for _, setter := range setters {
		setter(&options)
	}
	exists, err := f.FileExists(filePath)
	switch {
	case err != nil:
		return fmt.Errorf("cannot write to file: %w", err)
	case exists: // verify it is not a directory
		isDir, err := f.IsDirectory(filePath)
		if err != nil {
			return fmt.Errorf("cannot write to file: %w", err)
		} else if isDir {
			return fmt.Errorf("cannot write to file: %q is a directory", filePath)
		}
	case !exists: // create directories recursively
		parentDir := f.filepathDir(filePath)
		if err := f.mkdirAll(parentDir, 0700); err != nil {
			return fmt.Errorf("cannot write to file %q: %w", filePath, err)
		}
	}
	if err := f.writeFile(filePath, data, options.permissions); err != nil {
		return fmt.Errorf("cannot write to file %q: %w", filePath, err)
	}
	if options.ownership != nil {
		if err := f.chown(filePath, options.ownership.UID, options.ownership.GID); err != nil {
			return fmt.Errorf("cannot change ownership of file %q: %w", filePath, err)
		}
	}
	return nil
}
