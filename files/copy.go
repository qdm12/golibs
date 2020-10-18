package files

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
)

// CopyDirectory copies all files, directories and symlinks recursively to another path.
func (f *fileManager) CopyDirectory(fromPath, toPath string) error {
	errPrefix := fmt.Sprintf("cannot copy directory from %q to %q", fromPath, toPath)
	entries, err := f.readDir(fromPath)
	if err != nil {
		return fmt.Errorf("%s: %w", errPrefix, err)
	}
	for _, entry := range entries {
		subFromPath := filepath.Join(fromPath, entry.Name())
		subToPath := filepath.Join(toPath, entry.Name())
		fileInfo, err := os.Stat(subFromPath)
		if err != nil {
			return fmt.Errorf("%s: %w", errPrefix, err)
		}
		stat, ok := fileInfo.Sys().(*syscall.Stat_t)
		if !ok {
			return fmt.Errorf("%s: cannot get stats information for %q", errPrefix, subFromPath)
		}
		switch fileInfo.Mode() & os.ModeType { //nolint:exhaustive
		case os.ModeDir:
			const defaultPermissions os.FileMode = 0700
			if err := f.CreateDir(subToPath, Permissions(defaultPermissions)); err != nil {
				return fmt.Errorf("%s: %w", errPrefix, err)
			}
			if err := f.CopyDirectory(subFromPath, subToPath); err != nil {
				return fmt.Errorf("%s: %w", errPrefix, err)
			}
		case os.ModeSymlink:
			if err := f.CopySymLink(subFromPath, subToPath); err != nil {
				return fmt.Errorf("%s: %w", errPrefix, err)
			}
		default:
			if err := f.CopyFile(subFromPath, subToPath); err != nil {
				return fmt.Errorf("%s: %w", errPrefix, err)
			}
		}
		if err := os.Lchown(subToPath, int(stat.Uid), int(stat.Gid)); err != nil {
			return fmt.Errorf("%s: %w", errPrefix, err)
		}
		if isSymlink := entry.Mode()&os.ModeSymlink != 0; !isSymlink {
			if err := os.Chmod(subToPath, entry.Mode()); err != nil {
				return fmt.Errorf("%s: %w", errPrefix, err)
			}
		}
	}
	return nil
}

// CopyFile copies a file from a path to another path.
func (f *fileManager) CopyFile(fromPath, toPath string) (err error) {
	errPrefix := fmt.Sprintf("cannot copy file from %q to %q", fromPath, toPath)
	out, err := f.create(toPath)
	if err != nil {
		return fmt.Errorf("%s: %w", errPrefix, err)
	}

	defer func() {
		closeErr := out.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("%s: %w", errPrefix, closeErr)
		}
	}()

	in, err := f.open(fromPath)

	defer func() {
		closeErr := in.Close()
		if err == nil && closeErr != nil {
			err = fmt.Errorf("%s: %w", errPrefix, closeErr)
		}
	}()

	if err != nil {
		return fmt.Errorf("%s: %w", errPrefix, err)
	}
	if _, err := io.Copy(out, in); err != nil {
		return fmt.Errorf("%s: %w", errPrefix, err)
	}
	return nil
}

// CopySymLink copies a symlink to another path.
func (f *fileManager) CopySymLink(fromPath, toPath string) error {
	link, err := f.readlink(fromPath)
	if err != nil {
		return fmt.Errorf("cannot copy symlink from %q to %q: %w", fromPath, toPath, err)
	}
	return f.symlink(link, toPath)
}
