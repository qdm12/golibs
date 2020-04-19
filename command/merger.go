package command

import (
	"bufio"
	"context"
	"fmt"
	"io"

	"github.com/fatih/color"
)

// StreamMerger contains methods to merge some IO streams and process them line by line.
// It was designed to use for processing the outputs of multiple programs in a single goroutine.
//go:generate mockgen -destination=mock_command/streammerger.go . StreamMerger
type StreamMerger interface {
	// Merge merges the given stream with the given name to the CollectLines method, and should be run in a goroutine
	Merge(stream io.ReadCloser, setters ...MergeOptionSetter)
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

type mergeOptions struct {
	name  string
	color *color.Color
}

type MergeOptionSetter func(o *mergeOptions)

func MergeName(name string) MergeOptionSetter {
	return func(options *mergeOptions) {
		options.name = name
	}
}

func MergeColor(color color.Color) MergeOptionSetter {
	return func(options *mergeOptions) {
		options.color = &color
	}
}

// Merge merges the given stream with the given name to the CollectLines method, and should be run in a goroutine
func (s *streamMerger) Merge(stream io.ReadCloser, setters ...MergeOptionSetter) {
	options := &mergeOptions{}
	for _, setter := range setters {
		setter(options)
	}
	prefix := ""
	if len(options.name) > 0 {
		prefix = options.name + ": "
	}
	lineWrapper := func(s string) string {
		return prefix + s
	}
	if options.color != nil {
		lineWrapper = func(s string) string {
			return options.color.Sprintf(prefix + s)
		}
	}
	defer stream.Close()
	go func() { // Read lines infinitely
		scanner := bufio.NewScanner(stream)
		for scanner.Scan() {
			line := string(scanner.Bytes())
			s.chLine <- lineWrapper(line)
		}
		if err := scanner.Err(); err != nil {
			s.chErr <- fmt.Errorf("%sstream error: %w", prefix, err)
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
