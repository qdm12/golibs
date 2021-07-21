package files

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

var (
	ErrReadDirectory   = errors.New("cannot read directory")
	ErrStatFile        = errors.New("cannot stat file")
	ErrSyscallStat     = errors.New("file does not have syscall stat")
	ErrCreateDirectory = errors.New("cannot create directory")
	ErrCopyDirectory   = errors.New("cannot copy directory")
	ErrCopySymLink     = errors.New("cannot copy symlink")
	ErrCopyFile        = errors.New("cannot copy file")
	ErrChown           = errors.New("cannot change ownership")
	ErrChmod           = errors.New("cannot change permission mode")
)

// CopyDirectory copies all files, directories and symlinks recursively to another path.
func (f *fileManager) CopyDirectory(fromPath, toPath string) error {
	entries, err := f.readDir(fromPath)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrReadDirectory, err)
	}
	for _, entry := range entries {
		subFromPath := filepath.Join(fromPath, entry.Name())
		subToPath := filepath.Join(toPath, entry.Name())
		fileInfo, err := os.Stat(subFromPath)
		if err != nil {
			return fmt.Errorf("%w: for path %s: %s",
				ErrStatFile, subFromPath, err)
		}
		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("%w: for path %s", ErrSyscallStat, subFromPath)
		}

		mode := fileInfo.Mode() & os.ModeType
		switch mode {
		case os.ModeDir:
			const defaultPermissions os.FileMode = 0700
			if err := f.CreateDir(subToPath, Permissions(defaultPermissions)); err != nil {
				return fmt.Errorf("%w: %s", ErrCreateDirectory, err)
			}
			if err := f.CopyDirectory(subFromPath, subToPath); err != nil {
				return fmt.Errorf("%w: from %s to %s: %s",
					ErrCopyDirectory, subFromPath, subToPath, err)
			}
		case os.ModeSymlink:
			if err := f.CopySymLink(subFromPath, subToPath); err != nil {
				return fmt.Errorf("%w: from %s to %s: %s",
					ErrCopySymLink, subFromPath, subToPath, err)
			}
		default:
			if err := f.CopyFile(subFromPath, subToPath); err != nil {
				return fmt.Errorf("%w: from %s to %s: %s",
					ErrCopyFile, subFromPath, subToPath, err)
			}
		}
		if err := os.Lchown(subToPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return fmt.Errorf("%w: path %s: %s", ErrChown, subToPath, err)
		}
		if isSymlink := entry.Mode()&os.ModeSymlink != 0; !isSymlink {
			if err := os.Chmod(subToPath, entry.Mode()); err != nil {
				return fmt.Errorf("%w: path %s: %s", ErrChmod, subToPath, err)
			}
		}
	}
	return nil
}

var (
	ErrCreateFile           = errors.New("cannot create file")
	ErrOpenFile             = errors.New("cannot open file")
	ErrCloseSourceFile      = errors.New("cannot close source file")
	ErrCloseDestinationFile = errors.New("cannot close destination file")
	ErrCopyData             = errors.New("cannot copy data")
)

// CopyFile copies a file from a path to another path.
func (f *fileManager) CopyFile(fromPath, toPath string) (err error) {
	out, err := f.create(toPath)
	if err != nil {
		return fmt.Errorf("%w: path %s: %s", ErrCreateFile, toPath, err)
	}

	defer func() {
		closeErr := out.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("%w: path %s: %s",
				ErrCloseDestinationFile, toPath, closeErr)
		}
	}()

	in, err := f.open(fromPath)
	if err != nil {
		return fmt.Errorf("%w: path %s: %s", ErrOpenFile, fromPath, err)
	}

	defer func() {
		closeErr := in.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("%w: path %s: %s",
				ErrCloseSourceFile, fromPath, closeErr)
		}
	}()

	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("%w: %s", ErrCopyData, err)
	}

	return nil
}

var (
	ErrReadSymlink = errors.New("cannot read symlink")
	ErrMakeSymlink = errors.New("cannot make symlink")
)

// CopySymLink copies a symlink to another path.
func (f *fileManager) CopySymLink(fromPath, toPath string) error {
	link, err := f.readlink(fromPath)
	if err != nil {
		return fmt.Errorf("%w: at path %s: %s", ErrReadSymlink, fromPath, err)
	}

	if err := f.symlink(link, toPath); err != nil {
		return fmt.Errorf("%w: at path %s: %s", ErrReadSymlink, toPath, err)
	}

	return nil
}
