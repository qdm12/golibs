package logging

import "github.com/fatih/color"

type Color func() *color.Color

func ColorFgCyan() *color.Color {
	return color.New(color.FgCyan)
}

func ColorFgHiMagenta() *color.Color {
	return color.New(color.FgHiMagenta)
}
