package command

import (
	"context"
	"os/exec"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Start(t *testing.T) {
	t.Parallel()
	commander := &commander{
		execCommand: exec.CommandContext,
	}
	ctx := context.Background()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	stdoutLines := make(chan string)
	stderrLines := make(chan string)
	wait := make(chan error)

	expectedStdoutLines := []string{"hello", "world"}
	expectedStderrLines := []string{}

	go commander.Start(ctx, wg, stdoutLines, stderrLines, wait, "echo", "hello\nworld")

	var stdoutLinesIndex, stderrLinesIndex int

	var done bool
	for !done {
		select {
		case line := <-stdoutLines:
			assert.Equal(t, expectedStdoutLines[stdoutLinesIndex], line)
			stdoutLinesIndex++
		case line := <-stderrLines:
			assert.Equal(t, expectedStderrLines[stderrLinesIndex], line)
			stderrLinesIndex++
		case err := <-wait:
			assert.NoError(t, err)
			close(stdoutLines)
			close(stderrLines)
			close(wait)
			done = true
		}
	}
	wg.Wait()
}
