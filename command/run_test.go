package command

import (
	"errors"
	"testing"

	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_commander_Run(t *testing.T) {
	t.Parallel()

	errDummy := errors.New("dummy")

	testCases := map[string]struct {
		stdout []byte
		cmdErr error
		output string
		err    error
	}{
		"no output": {},
		"cmd error": {
			stdout: []byte("'hello \nworld'\n"),
			cmdErr: errDummy,
			output: "hello \nworld",
			err:    errDummy,
		},
	}

	for name, testCase := range testCases {
		testCase := testCase
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)

			commander := &commander{}
			mockCmd := NewMockCmd(ctrl)

			mockCmd.EXPECT().CombinedOutput().Return(testCase.stdout, testCase.cmdErr)

			output, err := commander.Run(mockCmd)

			if testCase.err != nil {
				require.Error(t, err)
				assert.Equal(t, testCase.err.Error(), err.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, testCase.output, output)
		})
	}
}
