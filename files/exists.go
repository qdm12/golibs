package files

import (
	"os"
)

// FilepathExists returns true if a file path exists.
func (f *FileManager) FilepathExists(filePath string) (exists bool, err error) {
	_, err = os.Stat(filePath)
	if err == nil {
		return true, nil
	} else if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// FileExists returns true if a file exists at the path given.
// If a directory is at the path, it returns false.
func (f *FileManager) FileExists(filePath string) (exists bool, err error) {
	info, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	default:
		return !info.IsDir(), nil
	}
}

// DirectoryExists returns true if a directory exists.
func (f *FileManager) DirectoryExists(filePath string) (exists bool, err error) {
	info, err := os.Stat(filePath)
	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	default:
		return info.IsDir(), nil
	}
}

// IsFile returns true if the path points to a file.
func (f *FileManager) IsFile(filePath string) (bool, error) {
	isDir, err := f.IsDirectory(filePath)
	return !isDir, err
}

// IsDirectory returns true if the path points to a directory.
func (f *FileManager) IsDirectory(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
