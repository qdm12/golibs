package files

// Remove removes a file or directory.
func (f *FileManager) Remove(filePath string) (err error) {
	return f.rm(filePath)
}
