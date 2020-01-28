package command

import (
	"os/exec"
	"strings"
)

type Command interface {
	Run(name string, arg ...string) (output string, err error)
}

type command struct {
	execCommand func(name string, arg ...string) *exec.Cmd
}

func NewCommand() Command {
	return &command{
		execCommand: exec.Command,
	}
}

func (c *command) Run(name string, arg ...string) (output string, err error) {
	cmd := c.execCommand(name, arg...)
	stdout, err := cmd.Output()
	if err != nil {
		return "", err
	}
	output = string(stdout)
	output = strings.TrimSuffix(output, "\n")
	lines := stringToLines(output)
	for i := range lines {
		lines[i] = strings.TrimPrefix(lines[i], "'")
		lines[i] = strings.TrimSuffix(lines[i], "'")
	}
	output = strings.Join(lines, "\n")
	return output, nil
}

func stringToLines(s string) (lines []string) {
	s = strings.TrimSuffix(s, "\n")
	return strings.Split(s, "\n")
}
