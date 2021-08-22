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

// New returns a new Env object which will read
// environment variables using os.GetEnv.
func New() *Env {
	return &Env{
		getuid: os.Getuid,
		getenv: os.Getenv,
		unset:  os.Unsetenv,
		fpAbs:  filepath.Abs,
	}
}
