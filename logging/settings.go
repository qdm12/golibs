package logging

import (
	"io"
	"os"

	"github.com/fatih/color"
)

type settings struct {
	level      Level
	caller     Caller
	color      *color.Color
	prefix     string
	writer     io.Writer
	preProcess func(line string) string
}

func newSettings(setters ...Option) (s settings) {
	// Defaults
	s.level = LevelInfo
	s.writer = os.Stdout

	for _, setter := range setters {
		setter(&s)
	}

	return s
}
