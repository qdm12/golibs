package logging

import (
	"io"
)

type Option func(s *settings)

func SetLevel(level Level) Option {
	return func(s *settings) {
		s.level = level
	}
}

func SetCaller(caller Caller) Option {
	return func(s *settings) {
		s.caller = caller
	}
}

func SetColor(color Color) Option {
	return func(s *settings) {
		s.color = color()
	}
}

func SetPrefix(prefix string) Option {
	return func(s *settings) {
		s.prefix = prefix
	}
}

func SetWriter(writer io.Writer) Option {
	return func(s *settings) {
		s.writer = writer
	}
}

// SetPreProcess sets a pre processing function to act as a middleware
// on the line before it is written to stdout, allowing it to get modified
// depending on custom logic. This does not block other child loggers
// but blocks the current logger, so it should be designed to be fast enough.
// Also note that the line is formatted and colored, but the prefix is not
// part of it.
func SetPreProcess(preProcess func(line string) string) Option {
	return func(s *settings) {
		s.preProcess = preProcess
	}
}
