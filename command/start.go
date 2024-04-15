package command

import (
	"bufio"
	"errors"
	"io"
	"os"
	"sync"
)

// Start launches a command and streams stdout and stderr to channels.
// All the channels returned should be closed when an error,
// nil or not, is received in the waitError channel.
// The channels should NOT be closed if an error is returned directly
// with err, as they will already be closed internally by the function.
func (c *Cmder) Start(cmd ExecCmd) (
	stdoutLines, stderrLines chan string, waitError chan error, err error) {
	wg := &sync.WaitGroup{}

	done := make(chan struct{})

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	wg.Add(1)
	stdoutLines = make(chan string)
	go streamToChannel(done, wg, stdout, stdoutLines)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		_ = stdout.Close()
		close(stdoutLines)
		wg.Wait()
		return nil, nil, nil, err
	}
	wg.Add(1)
	stderrLines = make(chan string)
	go streamToChannel(done, wg, stderr, stderrLines)

	err = cmd.Start()
	if err != nil {
		close(done)
		_ = stdout.Close()
		close(stdoutLines)
		_ = stderr.Close()
		close(stderrLines)
		wg.Wait()
		return nil, nil, nil, err
	}

	waitError = make(chan error)
	go func() {
		err := cmd.Wait()
		close(done)
		_ = stdout.Close()
		_ = stderr.Close()
		wg.Wait()
		waitError <- err
	}()

	return stdoutLines, stderrLines, waitError, nil
}

func streamToChannel(done <-chan struct{}, wg *sync.WaitGroup,
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
	if err == nil {
		return
	}

	// ignore the error if we are done
	select {
	case <-done:
		return
	default:
	}

	if !errors.Is(err, os.ErrClosed) {
		lines <- "stream error: " + err.Error()
	}
}
