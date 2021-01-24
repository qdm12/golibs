package command

import (
	"context"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Start(t *testing.T) {
	t.Parallel()
	commander := &commander{
		execCommand: exec.CommandContext,
	}
	ctx := context.Background()

	const program = "echo"
	const args = "hello\nworld"
	expectedStdoutLines := []string{"hello", "world"}
	expectedStderrLines := []string{}

	stdoutLines, stderrLines, waitError, err := commander.Start(ctx, program, args)

	assert.NoError(t, err)
	var stdoutLinesIndex, stderrLinesIndex int

	for {
		select {
		case line := <-stdoutLines:
			assert.Equal(t, expectedStdoutLines[stdoutLinesIndex], line)
			stdoutLinesIndex++
		case line := <-stderrLines:
			assert.Equal(t, expectedStderrLines[stderrLinesIndex], line)
			stderrLinesIndex++
		case err := <-waitError:
			assert.NoError(t, err)
			close(stdoutLines)
			close(stderrLines)
			close(waitError)
			return
		}
	}
}
