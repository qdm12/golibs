package files

func (f *FileManager) GetOwnership(filePath string) (userID, groupID int, err error) {
	panic("ownership not supported on Windows")
}
