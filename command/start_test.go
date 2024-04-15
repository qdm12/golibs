package command

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func linesToReadCloser(lines []string) io.ReadCloser {
	s := strings.Join(lines, "\n")
	return io.NopCloser(bytes.NewBufferString(s))
}

func Test_commander_Start(t *testing.T) {
	t.Parallel()

	errDummy := errors.New("dummy")

	testCases := map[string]struct {
		makeExecCmd func(ctrl *gomock.Controller) *MockExecCmd
		stdout      []string
		stderr      []string
		startErr    error
		waitErr     error
		err         error
	}{
		"no_output": {
			makeExecCmd: func(ctrl *gomock.Controller) *MockExecCmd {
				cmd := NewMockExecCmd(ctrl)
				cmd.EXPECT().StdoutPipe().Return(linesToReadCloser(nil), nil)
				cmd.EXPECT().StderrPipe().Return(linesToReadCloser(nil), nil)
				cmd.EXPECT().Start().Return(nil)
				cmd.EXPECT().Wait().Return(nil)
				return cmd
			},
		},
		"success": {
			makeExecCmd: func(ctrl *gomock.Controller) *MockExecCmd {
				cmd := NewMockExecCmd(ctrl)
				cmd.EXPECT().StdoutPipe().
					Return(linesToReadCloser([]string{"hello", "world"}), nil)
				cmd.EXPECT().StderrPipe().
					Return(linesToReadCloser([]string{"some", "error"}), nil)
				cmd.EXPECT().Start().Return(nil)
				cmd.EXPECT().Wait().Return(nil)
				return cmd
			},
			stdout: []string{"hello", "world"},
			stderr: []string{"some", "error"},
		},
		"stdout_pipe_error": {
			makeExecCmd: func(ctrl *gomock.Controller) *MockExecCmd {
				cmd := NewMockExecCmd(ctrl)
				cmd.EXPECT().StdoutPipe().Return(nil, errDummy)
				return cmd
			},
			err: errDummy,
		},
		"stderr_pipe_error": {
			makeExecCmd: func(ctrl *gomock.Controller) *MockExecCmd {
				cmd := NewMockExecCmd(ctrl)
				cmd.EXPECT().StdoutPipe().Return(linesToReadCloser(nil), nil)
				cmd.EXPECT().StderrPipe().Return(nil, errDummy)
				return cmd
			},
			err: errDummy,
		},
		"start error": {
			makeExecCmd: func(ctrl *gomock.Controller) *MockExecCmd {
				cmd := NewMockExecCmd(ctrl)
				cmd.EXPECT().StdoutPipe().Return(linesToReadCloser(nil), nil)
				cmd.EXPECT().StderrPipe().Return(linesToReadCloser(nil), nil)
				cmd.EXPECT().Start().Return(errDummy)
				return cmd
			},
			startErr: errDummy,
			err:      errDummy,
		},
		"wait error": {
			makeExecCmd: func(ctrl *gomock.Controller) *MockExecCmd {
				cmd := NewMockExecCmd(ctrl)
				cmd.EXPECT().StdoutPipe().Return(linesToReadCloser(nil), nil)
				cmd.EXPECT().StderrPipe().Return(linesToReadCloser(nil), nil)
				cmd.EXPECT().Start().Return(nil)
				cmd.EXPECT().Wait().Return(errDummy)
				return cmd
			},
			waitErr: errDummy,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)

			cmder := &Cmder{}
			mockCmd := testCase.makeExecCmd(ctrl)

			stdoutLines, stderrLines, waitError, err := cmder.Start(mockCmd)

			assert.ErrorIs(t, err, testCase.err)
			if testCase.err != nil {
				assert.Nil(t, stdoutLines)
				assert.Nil(t, stderrLines)
				assert.Nil(t, waitError)
				return
			}

			require.NoError(t, err)

			var stdoutIndex, stderrIndex int
			done := false
			for !done {
				select {
				case line := <-stdoutLines:
					assert.Equal(t, testCase.stdout[stdoutIndex], line)
					stdoutIndex++
				case line := <-stderrLines:
					assert.Equal(t, testCase.stderr[stderrIndex], line)
					stderrIndex++
				case err := <-waitError:
					assert.ErrorIs(t, err, testCase.waitErr)
					done = true
				}
			}

			assert.Equal(t, len(testCase.stdout), stdoutIndex)
			assert.Equal(t, len(testCase.stderr), stderrIndex)
		})
	}
}
