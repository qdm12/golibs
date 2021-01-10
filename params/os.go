package params

import (
	"os"
	"path/filepath"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . OS

// OS has methods to obtain values from the OS.
type OS interface {
	// UID obtains the user ID running the program.
	UID() int
	// GID obtains the group ID of the user running the program.
	GID() int
	// ExeDir obtains the executable directory.
	ExeDir() (dir string, err error)
}

type osImpl struct {
	getuid     func() int
	getgid     func() int
	executable func() (string, error)
}

func NewOS() OS {
	return &osImpl{
		getuid:     os.Getuid,
		getgid:     os.Getgid,
		executable: os.Executable,
	}
}

// UID obtains the user ID running the program.
func (o *osImpl) UID() int {
	return o.getuid()
}

// GID obtains the group ID of the user running the program.
func (o *osImpl) GID() int {
	return o.getgid()
}

// ExeDir obtains the executable directory.
func (o *osImpl) ExeDir() (dir string, err error) {
	ex, err := o.executable()
	if err != nil {
		return dir, err
	}
	dir = filepath.Dir(ex)
	return dir, nil
}
