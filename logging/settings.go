package logging

import (
	"io"
	"os"
)

type Settings struct {
	Level  Level
	Caller Caller
	Color  Color
	Prefix string
	Writer io.Writer
	// PreProcess sets a pre processing function to act as a middleware
	// on the line before it is written to stdout, allowing it to get modified
	// depending on custom logic. This does not block other child loggers
	// but blocks the current logger, so it should be designed to be fast enough.
	// Also note that the line is formatted and colored, but the prefix is not
	// part of it.
	PreProcess func(line string) string
}

// setDefaults only set the defaults where the value is empty.
func (s *Settings) setDefaults() {
	// Defaults
	if s.Writer == nil {
		s.Writer = os.Stdout
	}
}

// setEmptyValuesWith set each field value with the corresponding value
// in the other Settings if the field is set to its zero value.
func (s *Settings) setEmptyValuesWith(other Settings) {
	if s.Level == 0 {
		s.Level = other.Level
	}
	if s.Caller == 0 {
		s.Caller = other.Caller
	}
	if s.Prefix == "" {
		s.Prefix = other.Prefix
	}
	if s.Writer == nil {
		s.Writer = other.Writer
	}
	if s.PreProcess == nil {
		s.PreProcess = other.PreProcess
	}
}
