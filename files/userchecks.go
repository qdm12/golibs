package files

import (
	"fmt"
	"io/fs"
	"os"
	"syscall"
)

type accessible string

const (
	readable   accessible = "readable"
	writable   accessible = "writable"
	executable accessible = "executable"
)

func (a accessible) toBitMode() (bit fs.FileMode) {
	switch a {
	case readable:
		return 4 //nolint:gomnd
	case writable:
		return 2 //nolint:gomnd
	case executable:
		return 1
	default:
		return 0
	}
}

func IsReadable(filePath string, uid, gid int) (bool, error) {
	return isAccessible(filePath, uid, gid, readable)
}

func IsWritable(filePath string, uid, gid int) (bool, error) {
	return isAccessible(filePath, uid, gid, writable)
}

func IsExecutable(filePath string, uid, gid int) (bool, error) {
	return isAccessible(filePath, uid, gid, executable)
}

func isAccessible(filePath string, uid, gid int, accessibility accessible) (
	ok bool, err error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return false, err
	}
	mode := info.Mode()
	perm := mode.Perm()
	stat, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		panic(fmt.Sprintf("file %s does not have syscall stat", filePath))
	}

	relevantBit := accessibility.toBitMode()

	// Others
	if perm&relevantBit != 0 {
		return true, nil
	}

	// Group
	relevantBit *= 8
	if gid == int(stat.Gid) && perm&relevantBit != 0 {
		return true, nil
	}

	// User
	relevantBit *= 8
	if uid == int(stat.Uid) && perm&relevantBit != 0 {
		return true, nil
	}

	return false, nil
}
