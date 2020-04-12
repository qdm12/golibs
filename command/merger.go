package command

import (
	"bufio"
	"context"
	"fmt"
	"io"
)

//go:generate mockgen -destination=mock_command/streammerger.go . StreamMerger
type StreamMerger interface {
	// Merge merges the given stream with the given name to the CollectLines method, and should be run in a goroutine
	Merge(name string, stream io.ReadCloser)
	// CollectLines collects lines received from each stream merged in the streammerger and runs onNewLine on each line
	CollectLines(onNewLine func(line string)) error
}

type streamMerger struct {
	ctx    context.Context
	cancel context.CancelFunc
	chLine chan string
	chErr  chan error
}

func NewStreamMerger(ctx context.Context) StreamMerger {
	ctx, cancel := context.WithCancel(ctx)
	return &streamMerger{
		ctx:    ctx,
		cancel: cancel,
		chLine: make(chan string),
		chErr:  make(chan error),
	}
}

// Merge merges the given stream with the given name to the CollectLines method, and should be run in a goroutine
func (s *streamMerger) Merge(name string, stream io.ReadCloser) {
	defer stream.Close()
	go func() { // Read lines infinitely
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			line := string(scanner.Bytes())
			s.chLine <- fmt.Sprintf("%s: %s", name, line)
		}
		if err := scanner.Err(); err != nil {
			s.chErr <- fmt.Errorf("%s: stream error: %w", name, err)
		}
	}()
	<-s.ctx.Done() // blocks until context is canceled
}

// CollectLines collects lines received from each stream merged in the streammerger and runs onNewLine on each line
func (s *streamMerger) CollectLines(onNewLine func(line string)) error {
	defer func() {
		s.cancel() // stops other streams
		close(s.chLine)
		close(s.chErr)
	}()
	for {
		select {
		case line := <-s.chLine:
			onNewLine(line)
		case err := <-s.chErr:
			return err
		case <-s.ctx.Done():
			return nil
		}
	}
}
