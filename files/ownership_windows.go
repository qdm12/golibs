package files

func (f *FileManager) GetOwnership(filePath string) (userID, groupID int, err error) {
	panic("ownership not supported on Windows")
}

func (f *FileManager) SetOwnership(filePath string, userID, groupID int) error {
	panic("ownership not supported on Windows")
}
