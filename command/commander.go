package command

import (
	"io"
	"os/exec"
	"strings"
)

type Commander interface {
	Run(name string, arg ...string) (output string, err error)
	Start(name string, arg ...string) (stdoutPipe io.ReadCloser, waitFn func() error, err error)
}

type commander struct {
	execCommand func(name string, arg ...string) *exec.Cmd
}

func NewCommander() Commander {
	return &commander{
		execCommand: exec.Command,
	}
}

func (c *commander) Run(name string, arg ...string) (output string, err error) {
	cmd := c.execCommand(name, arg...)
	stdout, err := cmd.CombinedOutput()
	output = string(stdout)
	output = strings.TrimSuffix(output, "\n")
	lines := stringToLines(output)
	for i := range lines {
		lines[i] = strings.TrimPrefix(lines[i], "'")
		lines[i] = strings.TrimSuffix(lines[i], "'")
	}
	output = strings.Join(lines, "\n")
	return output, err
}

func (c *commander) Start(name string, arg ...string) (stdoutPipe io.ReadCloser, waitFn func() error, err error) {
	cmd := c.execCommand(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}
	return stdout, cmd.Wait, nil
}

func stringToLines(s string) (lines []string) {
	s = strings.TrimSuffix(s, "\n")
	return strings.Split(s, "\n")
}
