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
