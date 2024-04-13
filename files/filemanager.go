package files

import (
	"io/fs"
	"os"
	"path/filepath"
)

type FileManager struct {
	fileStat    func(name string) (os.FileInfo, error)
	isNotExist  func(err error) bool
	readFile    func(filename string) ([]byte, error)
	filepathDir func(path string) string
	mkdirAll    func(path string, perm os.FileMode) error
	writeFile   func(filename string, data []byte, perm os.FileMode) error
	chown       func(name string, uid int, gid int) error
	chmod       func(name string, mod os.FileMode) error
	rm          func(path string) error
	create      func(name string) (*os.File, error)
	open        func(name string) (*os.File, error)
	readlink    func(name string) (string, error)
	symlink     func(oldName, newName string) error
	readDir     func(dirname string) ([]fs.DirEntry, error)
}

func NewFileManager() *FileManager {
	return &FileManager{
		fileStat:    os.Stat,
		isNotExist:  os.IsNotExist,
		readFile:    os.ReadFile,
		filepathDir: filepath.Dir,
		mkdirAll:    os.MkdirAll,
		writeFile:   os.WriteFile,
		chown:       os.Chown,
		chmod:       os.Chmod,
		rm:          os.RemoveAll,
		create:      os.Create,
		open:        os.Open,
		readlink:    os.Readlink,
		symlink:     os.Symlink,
		readDir:     os.ReadDir,
	}
}
