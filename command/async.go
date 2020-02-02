package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

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

func (c *commander) MergeLineReaders(ctx context.Context, onNewLine func(line string), readers map[string]io.ReadCloser) error {
	buffers := make(map[string]*bufio.Reader, len(readers))
	for name, reader := range readers {
		buffers[name] = bufio.NewReader(reader)
	}
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chLine := make(chan string)
	chErr := make(chan error)
	// Read from all buffers asynchronously
	for name, buffer := range buffers {
		go func(name string, buffer *bufio.Reader) {
			go func() { // Read lines infinitely
				for {
					line, _, err := buffer.ReadLine()
					if err != nil {
						chErr <- err
						return
					}
					chLine <- fmt.Sprintf("%s: %s", name, line)
				}
			}()
			<-ctx.Done() // blocks until context is canceled
			buffer.Reset(nil)
		}(name, buffer)
	}
	// Collect lines from all buffers synchronously
	for {
		select {
		case line := <-chLine:
			onNewLine(line)
		case err := <-chErr:
			cancel()
			close(chLine)
			return err
		case <-ctx.Done():
			close(chLine)
			return nil
		}
	}
}
