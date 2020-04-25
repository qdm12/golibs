package files

import "fmt"

type accessible string

const (
	readable   accessible = "readable"
	writable   accessible = "writable"
	executable accessible = "executable"
)

func (f *fileManager) IsReadable(filePath string, uid, gid int) (bool, error) {
	return f.isAccessible(filePath, uid, gid, readable)
}

func (f *fileManager) IsWritable(filePath string, uid, gid int) (bool, error) {
	return f.isAccessible(filePath, uid, gid, writable)
}

func (f *fileManager) IsExecutable(filePath string, uid, gid int) (bool, error) {
	return f.isAccessible(filePath, uid, gid, executable)
}

func (f *fileManager) isAccessible(filePath string, uid, gid int, accessibility accessible) (accessible bool, err error) {
	errPrefix := fmt.Sprintf("%s is not %s for user with uid %d and gid %d", filePath, accessibility, uid, gid)
	ownerUID, ownerGID, err := f.GetOwnership(filePath)
	if err != nil {
		return false, fmt.Errorf("%s: %w", errPrefix, err)
	}
	accessible = false
	switch accessibility {
	case readable:
		accessible, _, _, err = f.GetOthersPermissions(filePath)
	case writable:
		_, accessible, _, err = f.GetOthersPermissions(filePath)
	case executable:
		_, _, accessible, err = f.GetOthersPermissions(filePath)
	}
	if err != nil {
		return false, fmt.Errorf("%s: %w", errPrefix, err)
	} else if accessible {
		return true, nil
	}
	if gid == ownerGID {
		accessible := false
		switch accessibility {
		case readable:
			accessible, _, _, err = f.GetGroupPermissions(filePath)
		case writable:
			_, accessible, _, err = f.GetGroupPermissions(filePath)
		case executable:
			_, _, accessible, err = f.GetGroupPermissions(filePath)
		}
		if err != nil {
			return false, fmt.Errorf("%s: %w", errPrefix, err)
		} else if accessible {
			return true, nil
		}
	}
	if uid == ownerUID {
		accessible := false
		switch accessibility {
		case readable:
			accessible, _, _, err = f.GetGroupPermissions(filePath)
		case writable:
			_, accessible, _, err = f.GetGroupPermissions(filePath)
		case executable:
			_, _, accessible, err = f.GetGroupPermissions(filePath)
		}
		if err != nil {
			return false, fmt.Errorf("%s: %w", errPrefix, err)
		} else if accessible {
			return true, nil
		}
	}
	return false, nil
}
