package command

import (
	"context"
	"io"
	"os/exec"
	"strings"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Commander

// Commander contains methods to run and start shell commands
type Commander interface {
	// Run runs a command in a blocking manner, returning its output and an error if it failed
	Run(ctx context.Context, name string, arg ...string) (output string, err error)
	// Start launches a command asynchronously and returns streams for stdout, stderr as well as a wait function
	Start(ctx context.Context, name string, arg ...string) (stdoutPipe, stderrPipe io.ReadCloser, waitFn func() error, err error)
}

type commander struct {
	execCommand func(ctx context.Context, name string, arg ...string) *exec.Cmd
}

func NewCommander() Commander {
	return &commander{
		execCommand: exec.CommandContext,
	}
}

// Run runs a command in a blocking manner, returning its output and an error if it failed
func (c *commander) Run(ctx context.Context, name string, arg ...string) (output string, err error) {
	cmd := c.execCommand(ctx, name, arg...)
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

// Start launches a command asynchronously and returns streams for stdout, stderr as well as a wait function
func (c *commander) Start(ctx context.Context, name string, arg ...string) (stdoutPipe, stderrPipe io.ReadCloser, waitFn func() error, err error) {
	cmd := c.execCommand(ctx, name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, nil, nil, err
	}
	return stdout, stderr, cmd.Wait, nil
}

func stringToLines(s string) (lines []string) {
	s = strings.TrimSuffix(s, "\n")
	return strings.Split(s, "\n")
}
