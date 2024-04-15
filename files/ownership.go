package files

import (
	"os"
	"syscall"
)

// GetOwnership obtains the user ID and group ID owning the file or directory
// or returns 0 and 0 if running in Windows (no IDs).
func (f *FileManager) GetOwnership(filePath string) (userID, groupID int, err error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, 0, err
	}
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		panic("file does not have syscall stat")
	}
	return int(stat.Uid), int(stat.Gid), nil
}
