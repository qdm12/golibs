package params

import (
	"os"
	"path/filepath"
)

type Env struct {
	getuid func() int
	getenv func(key string) string
	unset  func(k string) error
	fpAbs  func(s string) (string, error)
}

// NewEnv returns a new Env object.
func NewEnv() *Env {
	return &Env{
		getuid: os.Getuid,
		getenv: os.Getenv,
		unset:  os.Unsetenv,
		fpAbs:  filepath.Abs,
	}
}
