package files

import "os"

type writeOptions struct {
	ownership *struct {
		UID int
		GID int
	}
	permissions os.FileMode
}

func newWriteOptions(defaultPermissions os.FileMode) writeOptions {
	return writeOptions{
		permissions: defaultPermissions,
	}
}

type WriteOptionSetter func(o *writeOptions)

func Ownership(uid, gid int) WriteOptionSetter {
	return func(options *writeOptions) {
		options.ownership = &struct {
			UID int
			GID int
		}{uid, gid}
	}
}

func Permissions(permissions os.FileMode) WriteOptionSetter {
	return func(options *writeOptions) {
		options.permissions = permissions
	}
}
