package files

import (
	"fmt"
)

// ReadFile reads the entire data of a file.
func (f *fileManager) ReadFile(filePath string) (data []byte, err error) {
	exists, err := f.FileExists(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file: %w", err)
	} else if !exists {
		return nil, fmt.Errorf("cannot read file: %q does not exist", filePath)
	}
	return f.readFile(filePath)
}
