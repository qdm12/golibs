package command

import (
	"os/exec"
	"strings"
)

func Run(command string, arg ...string) (output string, err error) {
	cmd := exec.Command(command, arg...)
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
