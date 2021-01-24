package command

import (
	"context"
	"strings"
)

// Run runs a command in a blocking manner, returning its output and an error if it failed.
func (c *commander) Run(ctx context.Context, name string, arg ...string) (output string, err error) {
	cmd := c.execCommand(ctx, name, arg...)
	stdout, err := cmd.CombinedOutput()
	output = string(stdout)
	output = strings.TrimSuffix(output, "\n")
	lines := stringToLines(output)
	for i := range lines {
		lines[i] = strings.TrimPrefix(lines[i], "'")
		lines[i] = strings.TrimSuffix(lines[i], "'")
	}
	output = strings.Join(lines, "\n")
	return output, err
}

func stringToLines(s string) (lines []string) {
	s = strings.TrimSuffix(s, "\n")
	return strings.Split(s, "\n")
}
