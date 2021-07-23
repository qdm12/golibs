package command

import "io"

//go:generate mockgen -destination=execcmd_mock_test.go -package=command . ExecCmd

// ExecCmd is the interface for exec.Cmd.
type ExecCmd interface {
	CombinedOutput() ([]byte, error)
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
}
