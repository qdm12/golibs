package params

import (
	"os"
	"path/filepath"
)

type Env struct {
	kv     map[string]string
	getuid func() int
	unset  func(k string) error
	fpAbs  func(s string) (string, error)
}

// New returns a new Env object which will read
// environment variables using os.GetEnv.
func New() *Env {
	return &Env{
		kv:     make(map[string]string),
		getuid: os.Getuid,
		unset:  os.Unsetenv,
		fpAbs:  filepath.Abs,
	}
}
