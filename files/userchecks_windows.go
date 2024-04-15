package files

func (f *FileManager) IsReadable(filePath string, uid, gid int) (bool, error) {
	return true, nil
}

func (f *FileManager) IsWritable(filePath string, uid, gid int) (bool, error) {
	return true, nil
}

func (f *FileManager) IsExecutable(filePath string, uid, gid int) (bool, error) {
	return true, nil
}
