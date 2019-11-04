package files

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// FileExists returns true if a file exists.
func FileExists(filePath string) (exists bool, err error) {
	_, err = os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// WriteLinesToFile writes a slice of strings as lines to a file.
// It creates any directory not existing in the file path if necessary.
func WriteLinesToFile(filePath string, lines []string) error {
	data := []byte(strings.Join(lines, "\n"))
	return WriteToFile(filePath, data)
}

// WriteToFile writes data bytes to a file, and creates any
// directory not existing in the file path if necessary.
func WriteToFile(filePath string, data []byte) error {
	exists, err := FileExists(filePath)
	if err != nil {
		return fmt.Errorf("cannot write to file %q: %w", filePath, err)
	}
	if !exists {
		parentDir := filepath.Dir(filePath)
		err := os.MkdirAll(parentDir, 0700)
		if err != nil {
			return fmt.Errorf("cannot write to file %q: %w", filePath, err)
		}
	}
	err = ioutil.WriteFile(filePath, data, 0700)
	return err
}
