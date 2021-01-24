package command

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"sync"
)

// Start launches a command and reads from its stdout and stderr streams
// until it completes. It should therefore be run in a goroutine.
// All the channels given should also be closed after an error,
// nil or not, is received in the wait channel.
func (c *commander) Start(ctx context.Context, wg *sync.WaitGroup,
	stdoutLines, stderrLines chan<- string, wait chan<- error,
	name string, arg ...string) {
	defer wg.Done()

	streamWg := &sync.WaitGroup{}
	defer streamWg.Wait()

	cmd := c.execCommand(ctx, name, arg...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		_ = stdout.Close()
		wait <- err
		return
	}
	streamWg.Add(1)
	go streamToChannel(ctx, streamWg, stdout, stdoutLines)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = stdout.Close()
		_ = stderr.Close()
		wait <- err
		return
	}
	streamWg.Add(1)
	go streamToChannel(ctx, streamWg, stdout, stderrLines)

	if err := cmd.Start(); err != nil {
		_ = stdout.Close()
		_ = stderr.Close()
		wait <- err
		return
	}

	wait <- cmd.Wait()
	_ = stdout.Close()
	_ = stderr.Close()
}

func streamToChannel(ctx context.Context, wg *sync.WaitGroup,
	stream io.Reader, lines chan<- string) {
	defer wg.Done()
	scanner := bufio.NewScanner(stream)
	for scanner.Scan() {
		// scanner is closed if the context is canceled
		// or if the command failed starting because the
		// stream is closed (io.EOF error).
		lines <- scanner.Text()
	}
	err := scanner.Err()
	if ctx.Err() == nil && err != nil && !errors.Is(err, os.ErrClosed) {
		lines <- "stream error: " + err.Error()
	}
}
