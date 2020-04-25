package files

// FilepathExists returns true if a fiel path exists.
func (f *fileManager) FilepathExists(filePath string) (exists bool, err error) {
	_, err = f.fileStat(filePath)
	if err == nil {
		return true, nil
	} else if f.isNotExist(err) {
		return false, nil
	}
	return false, err
}

// FileExists returns true if a file exists.
func (f *fileManager) FileExists(filePath string) (exists bool, err error) {
	exists, err = f.FilepathExists(filePath)
	if err != nil {
		return false, err
	} else if !exists {
		return false, nil
	}
	return f.IsFile(filePath)
}

// DirectoryExists returns true if a directory exists.
func (f *fileManager) DirectoryExists(filePath string) (exists bool, err error) {
	exists, err = f.FilepathExists(filePath)
	if err != nil {
		return false, err
	} else if !exists {
		return false, nil
	}
	return f.IsDirectory(filePath)
}

// IsFile returns true if the path points to a file
func (f *fileManager) IsFile(filePath string) (bool, error) {
	isDir, err := f.IsDirectory(filePath)
	return !isDir, err
}

// IsDirectory returns true if the path points to a directory
func (f *fileManager) IsDirectory(filePath string) (bool, error) {
	info, err := f.fileStat(filePath)
	if err != nil {
		return false, err
	}
	return info.IsDir(), nil
}
