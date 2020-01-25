package files

import (
	"fmt"
	"strings"
)

// WriteLinesToFile writes a slice of strings as lines to a file.
// It creates any directory not existing in the file path if necessary.
func (f *fileManager) WriteLinesToFile(filePath string, lines []string) error {
	data := []byte(strings.Join(lines, "\n"))
	return f.WriteToFile(filePath, data)
}

// Touch creates an empty file at the file path given.
// It creates any directory not existing in the file path if necessary.
func (f *fileManager) Touch(filePath string) error {
	return f.WriteToFile(filePath, nil)
}

// WriteToFile writes data bytes to a file, and creates any
// directory not existing in the file path if necessary.
func (f *fileManager) WriteToFile(filePath string, data []byte) error {
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
	return f.writeFile(filePath, data, 0600)
}
