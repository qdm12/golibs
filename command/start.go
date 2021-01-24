package command

import (
	"bufio"
	"context"
	"errors"
	"io"
	"os"
	"sync"
)

// Start launches a command and stream stdout and stderr to channels.
// All the channels returned should be closed when an error,
// nil or not, is received in the waitError channel.
func (c *commander) Start(ctx context.Context, name string, arg ...string) (
	stdoutLines, stderrLines chan string, waitError chan error, err error) {
	cmd := c.execCommand(ctx, name, arg...)

	wg := &sync.WaitGroup{}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	wg.Add(1)
	stdoutLines = make(chan string)
	go streamToChannel(ctx, wg, stdout, stdoutLines)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = stdout.Close()
		close(stdoutLines)
		return nil, nil, nil, err
	}
	wg.Add(1)
	stderrLines = make(chan string)
	go streamToChannel(ctx, wg, stderr, stderrLines)

	if err := cmd.Start(); err != nil {
		_ = stdout.Close()
		close(stdoutLines)
		_ = stderr.Close()
		close(stderrLines)
		return nil, nil, nil, err
	}

	waitError = make(chan error)
	go func() {
		err := cmd.Wait()
		_ = stdout.Close()
		_ = stderr.Close()
		wg.Wait()
		waitError <- err
	}()

	return stdoutLines, stderrLines, waitError, nil
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
