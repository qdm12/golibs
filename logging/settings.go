package logging

import (
	"io"
	"os"

	"github.com/fatih/color"
)

type settings struct {
	level       Level
	caller      Caller
	color       *color.Color
	prefix      string
	writer      io.Writer
	postProcess func(line string) string
}

func newSettings(setters ...Option) (s settings) {
	// Defaults
	s.level = LevelInfo
	s.writer = os.Stdout
	s.postProcess = func(line string) string { return line }

	for _, setter := range setters {
		setter(&s)
	}

	return s
}
