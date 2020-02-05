package files

import (
	"fmt"
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

// WriteToFile writes data bytes to a file, and creates any
// directory not existing in the file path if necessary.
func (f *fileManager) WriteToFile(filePath string, data []byte, setters ...WriteOptionSetter) error {
	options := newWriteOptions(0600)
	for _, setter := range setters {
		setter(&options)
	}
	exists, err := f.FileExists(filePath)
	if err != nil {
		return fmt.Errorf("cannot write to file: %w", err)
	} else if !exists {
		parentDir := f.filepathDir(filePath)
		err := f.mkdirAll(parentDir, 0700)
		if err != nil {
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
