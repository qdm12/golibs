package files

import (
	"errors"
)

var (
	ErrFileNotExist = errors.New("file does not exist")
)

// ReadFile reads the entire data of a file.
func (f *FileManager) ReadFile(filePath string) (data []byte, err error) {
	exists, err := f.FileExists(filePath)
	if err != nil {
		return nil, err
	} else if !exists {
		return nil, ErrFileNotExist
	}
	return f.readFile(filePath)
}
