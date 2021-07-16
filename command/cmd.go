package command

import "io"

//go:generate mockgen -destination=cmd_mock_test.go -package=command . Cmd

// Cmd is the interface for exec.Cmd.
type Cmd interface {
	CombinedOutput() ([]byte, error)
	StdoutPipe() (io.ReadCloser, error)
	StderrPipe() (io.ReadCloser, error)
	Start() error
	Wait() error
}
