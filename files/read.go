package files

import (
	"fmt"
	"io/fs"
	"os"
)

// ReadFile reads the entire data of a file.
func (f *FileManager) ReadFile(filePath string) (data []byte, err error) {
	exists, err := f.FileExists(filePath)
	if err != nil {
		return nil, fmt.Errorf("checking file existence: %w", err)
	} else if !exists {
		return nil, fmt.Errorf("%w: %s", fs.ErrNotExist, filePath)
	}
	return os.ReadFile(filePath)
}
