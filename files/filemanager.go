package files

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileManager interface {
	FilepathExists(filePath string) (exists bool, err error)
	FileExists(filePath string) (exists bool, err error)
	DirectoryExists(filePath string) (exists bool, err error)
}

type fileManager struct {
	fileStat    func(name string) (os.FileInfo, error)
	isNotExist  func(err error) bool
	readFile    func(filename string) ([]byte, error)
	filepathDir func(path string) string
	mkdirAll    func(path string, perm os.FileMode) error
	writeFile   func(filename string, data []byte, perm os.FileMode) error
}

func NewFileManager() FileManager {
	return &fileManager{
		fileStat:    os.Stat,
		isNotExist:  os.IsNotExist,
		readFile:    ioutil.ReadFile,
		filepathDir: filepath.Dir,
		mkdirAll:    os.MkdirAll,
		writeFile:   ioutil.WriteFile,
	}
}
