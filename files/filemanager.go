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
	GetOwnership(filePath string) (userID, groupID int, err error)
	GetUserPermissions(filePath string) (read, write, execute bool, err error)
	SetUserPermissions(filepath string, mod os.FileMode) error
	ReadFile(filePath string) (data []byte, err error)
	WriteLinesToFile(filePath string, lines []string, setters ...WriteOptionSetter) error
	Touch(filePath string, setters ...WriteOptionSetter) error
	WriteToFile(filePath string, data []byte, setters ...WriteOptionSetter) error
}

type fileManager struct {
	fileStat    func(name string) (os.FileInfo, error)
	isNotExist  func(err error) bool
	readFile    func(filename string) ([]byte, error)
	filepathDir func(path string) string
	mkdirAll    func(path string, perm os.FileMode) error
	writeFile   func(filename string, data []byte, perm os.FileMode) error
	chown       func(name string, uid int, gid int) error
	chmod       func(name string, mod os.FileMode) error
}

func NewFileManager() FileManager {
	return &fileManager{
		fileStat:    os.Stat,
		isNotExist:  os.IsNotExist,
		readFile:    ioutil.ReadFile,
		filepathDir: filepath.Dir,
		mkdirAll:    os.MkdirAll,
		writeFile:   ioutil.WriteFile,
		chown:       os.Chown,
		chmod:       os.Chmod,
	}
}
