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
	info, _ := f.fileStat(filePath)
	if info.IsDir() {
		return false, nil
	}
	return true, nil
}

// DirectoryExists returns true if a directory exists.
func (f *fileManager) DirectoryExists(filePath string) (exists bool, err error) {
	exists, err = f.FilepathExists(filePath)
	if err != nil {
		return false, err
	} else if !exists {
		return false, nil
	}
	info, _ := f.fileStat(filePath)
	if !info.IsDir() {
		return false, nil
	}
	return true, nil
}
