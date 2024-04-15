package files

import "os"

// Remove removes a file or directory.
func (f *FileManager) Remove(filePath string) (err error) {
	return os.Remove(filePath)
}
