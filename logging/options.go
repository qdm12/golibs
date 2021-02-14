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

func SetPostProcess(postProcess func(line string) string) Option {
	return func(s *settings) {
		s.postProcess = postProcess
	}
}
