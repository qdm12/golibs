package files

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

var (
	ErrReadDirectory   = errors.New("cannot read directory")
	ErrStatFile        = errors.New("cannot stat file")
	ErrCreateDirectory = errors.New("cannot create directory")
	ErrCopyDirectory   = errors.New("cannot copy directory")
	ErrCopySymLink     = errors.New("cannot copy symlink")
	ErrCopyFile        = errors.New("cannot copy file")
	ErrChown           = errors.New("cannot change ownership")
	ErrChmod           = errors.New("cannot change permission mode")
)

// CopyDirectory copies all files, directories and symlinks recursively to another path.
func (f *FileManager) CopyDirectory(fromPath, toPath string) error {
	entries, err := f.readDir(fromPath)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrReadDirectory, err)
	}
	for _, entry := range entries {
		err = f.copyDirEntry(fromPath, toPath, entry)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f *FileManager) copyDirEntry(fromPath, toPath string, entry fs.DirEntry) (err error) {
	subFromPath := filepath.Join(fromPath, entry.Name())
	subToPath := filepath.Join(toPath, entry.Name())
	fileInfo, err := entry.Info()
	if err != nil {
		return fmt.Errorf("%w: for path %s: %w",
			ErrStatFile, subFromPath, err)
	}

	switch {
	case fileInfo.Mode()&os.ModeDir != 0:
		const defaultPermissions os.FileMode = 0700
		err = f.CreateDir(subToPath, Permissions(defaultPermissions))
		if err != nil {
			return fmt.Errorf("%w: %w", ErrCreateDirectory, err)
		}
		err = f.CopyDirectory(subFromPath, subToPath)
		if err != nil {
			return fmt.Errorf("%w: from %s to %s: %w",
				ErrCopyDirectory, subFromPath, subToPath, err)
		}
	case fileInfo.Mode()&os.ModeSymlink != 0:
		err = f.CopySymLink(subFromPath, subToPath)
		if err != nil {
			return fmt.Errorf("%w: from %s to %s: %w",
				ErrCopySymLink, subFromPath, subToPath, err)
		}
	default:
		err = f.CopyFile(subFromPath, subToPath)
		if err != nil {
			return fmt.Errorf("%w: from %s to %s: %w",
				ErrCopyFile, subFromPath, subToPath, err)
		}
	}

	err = os.Chmod(subToPath, fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("%w: path %s: %w", ErrChmod, subToPath, err)
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
func (f *FileManager) CopyFile(fromPath, toPath string) (err error) {
	out, err := f.create(toPath)
	if err != nil {
		return fmt.Errorf("%w: path %s: %w", ErrCreateFile, toPath, err)
	}

	defer func() {
		closeErr := out.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("%w: path %s: %w",
				ErrCloseDestinationFile, toPath, closeErr)
		}
	}()

	in, err := f.open(fromPath)
	if err != nil {
		return fmt.Errorf("%w: path %s: %w", ErrOpenFile, fromPath, err)
	}

	defer func() {
		closeErr := in.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("%w: path %s: %w",
				ErrCloseSourceFile, fromPath, closeErr)
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCopyData, err)
	}

	return nil
}

var (
	ErrReadSymlink = errors.New("cannot read symlink")
	ErrMakeSymlink = errors.New("cannot make symlink")
)

// CopySymLink copies a symlink to another path.
func (f *FileManager) CopySymLink(fromPath, toPath string) error {
	link, err := f.readlink(fromPath)
	if err != nil {
		return fmt.Errorf("%w: at path %s: %w", ErrReadSymlink, fromPath, err)
	}

	err = f.symlink(link, toPath)
	if err != nil {
		return fmt.Errorf("%w: at path %s: %w", ErrReadSymlink, toPath, err)
	}

	return nil
}
