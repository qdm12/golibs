package command

import (
	"context"
	"os/exec"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Commander

// Commander contains methods to run and start shell commands.
type Commander interface {
	// Run runs a command in a blocking manner, returning its output and an error if it failed
	Run(ctx context.Context, name string, arg ...string) (output string, err error)
	// Start launches a command and stream stdout and stderr to channels.
	// All the channels returned should be closed when an error,
	// nil or not, is received in the waitError channel.
	Start(ctx context.Context, name string, arg ...string) (
		stdoutLines, stderrLines chan string, waitError chan error, err error)
}

type commander struct {
	execCommand func(ctx context.Context, name string, arg ...string) *exec.Cmd
}

func NewCommander() Commander {
	return &commander{
		execCommand: exec.CommandContext,
	}
}
