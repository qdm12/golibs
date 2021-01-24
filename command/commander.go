package command

import (
	"context"
	"os/exec"
	"sync"
)

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Commander

// Commander contains methods to run and start shell commands.
type Commander interface {
	// Run runs a command in a blocking manner, returning its output and an error if it failed
	Run(ctx context.Context, name string, arg ...string) (output string, err error)
	// Start launches a command and reads from its stdout and stderr streams
	// until it completes. It should therefore be run in a goroutine.
	// All the channels given should also be closed after an error,
	// nil or not, is received in the wait channel.
	Start(ctx context.Context, wg *sync.WaitGroup,
		stdoutLines, stderrLines chan<- string, wait chan<- error,
		name string, arg ...string)
}

type commander struct {
	execCommand func(ctx context.Context, name string, arg ...string) *exec.Cmd
}

func NewCommander() Commander {
	return &commander{
		execCommand: exec.CommandContext,
	}
}
