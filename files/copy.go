package files

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// CopyDirectory copies all files, directories and symlinks recursively to another path.
func CopyDirectory(fromPath, toPath string) error {
	entries, err := os.ReadDir(fromPath)
	if err != nil {
		return fmt.Errorf("reading directory: %w", err)
	}
	for _, entry := range entries {
		err = copyDirEntry(fromPath, toPath, entry)
		if err != nil {
			return fmt.Errorf("copying directory entry: %w", err)
		}
	}
	return nil
}

func copyDirEntry(fromPath, toPath string, entry fs.DirEntry) (err error) {
	subFromPath := filepath.Join(fromPath, entry.Name())
	subToPath := filepath.Join(toPath, entry.Name())
	fileInfo, err := entry.Info()
	if err != nil {
		return fmt.Errorf("stating file %s: %w", subFromPath, err)
	}

	switch {
	case fileInfo.Mode()&os.ModeDir != 0:
		const defaultPermissions os.FileMode = 0700
		err = CreateDir(subToPath, Permissions(defaultPermissions))
		if err != nil {
			return fmt.Errorf("creating directory: %w", err)
		}
		err = CopyDirectory(subFromPath, subToPath)
		if err != nil {
			return err // do not wrap due to the recursion
		}
	case fileInfo.Mode()&os.ModeSymlink != 0:
		err = CopySymLink(subFromPath, subToPath)
		if err != nil {
			return fmt.Errorf("copying symlink: %w", err)
		}
	default:
		err = CopyFile(subFromPath, subToPath)
		if err != nil {
			return fmt.Errorf("copying file: %w", err)
		}
	}

	err = os.Chmod(subToPath, fileInfo.Mode())
	if err != nil {
		return fmt.Errorf("changing permissions: %w", err)
	}
	return nil
}

// CopyFile copies a file from a path to another path.
func CopyFile(fromPath, toPath string) (err error) {
	out, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("creating file: %w", err)
	}

	defer func() {
		closeErr := out.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("closing destination file: %w", closeErr)
		}
	}()

	in, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("opening file: %w", err)
	}

	defer func() {
		closeErr := in.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("closing source file: %w", closeErr)
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return fmt.Errorf("copying data: %w", err)
	}

	return nil
}

// CopySymLink copies a symlink to another path.
func CopySymLink(fromPath, toPath string) error {
	link, err := os.Readlink(fromPath)
	if err != nil {
		return fmt.Errorf("reading source symlink: %w", err)
	}

	err = os.Symlink(link, toPath)
	if err != nil {
		return fmt.Errorf("creating desstination symlink: %w", err)
	}

	return nil
}
