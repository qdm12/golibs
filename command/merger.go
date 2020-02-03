package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"strings"
)

type StreamMerger interface {
	Add(name string, stream io.ReadCloser)
	ListenToAll(ctx context.Context, onNewLine func(line string)) error
}

type streamMerger struct {
	readers map[string]io.ReadCloser
}

func NewStreamMerger() StreamMerger {
	return &streamMerger{
		readers: make(map[string]io.ReadCloser),
	}
}

func (s *streamMerger) Add(name string, stream io.ReadCloser) {
	s.readers[name] = stream
}

func (s *streamMerger) ListenToAll(ctx context.Context, onNewLine func(line string)) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	chLine := make(chan string)
	chErr := make(chan error)
	// Read from all buffers asynchronously
	for name, reader := range s.readers {
		go func(name string, reader io.ReadCloser) {
			defer reader.Close()
			go func() { // Read lines infinitely
				scanner := bufio.NewScanner(reader)
				for scanner.Scan() {
					b := scanner.Bytes()
					for _, line := range strings.Split(string(b), "\n") {
						chLine <- fmt.Sprintf("%s: %s", name, line)
					}
				}
				if err := scanner.Err(); err != nil {
					chErr <- fmt.Errorf("%s: %w", name, err)
				}
			}()
			<-ctx.Done() // blocks until context is canceled
		}(name, reader)
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
