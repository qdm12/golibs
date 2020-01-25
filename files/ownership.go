package files

import (
	"fmt"
	"syscall"
)

// GetOwnership obtains the user ID and group ID owning the file or directory
// or returns 0 and 0 if running in Windows (no IDs)
func (f *fileManager) GetOwnership(filePath string) (userID, groupID int, err error) {
	info, err := f.fileStat(filePath)
	if err != nil {
		return 0, 0, fmt.Errorf("cannot get ownership: %w", err)
	}
	info.Mode()
	stat, ok := info.Sys().(*syscall.Stat_t) // TODO change to use stat directly
	if !ok {
		return 0, 0, nil // Windows
	}
	return int(stat.Uid), int(stat.Gid), nil
}
