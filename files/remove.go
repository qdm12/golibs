package files

// Remove removes a file or directory.
func (f *fileManager) Remove(filePath string) (err error) {
	return f.rm(filePath)
}
