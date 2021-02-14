package logging

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_stdLog_Debug(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBuffer(nil)

	logger := New(StdLog, SetWriter(buffer), SetCaller(CallerShort))

	logger.Debug("isn't this %q...", "function")
	logger.Debug("...fun?")

	result := buffer.String()
	buffer.Reset()

	result = strings.TrimSuffix(result, "\n")

	lines := strings.Split(result, "\n")
	require.Len(t, lines, 2)

	expectedVariablePrefix, err := regexp.Compile(`2[0-9]{3}/[0-1][0-9]/[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] `)
	require.NoError(t, err)

	expectedLinesWithoutPrefix := []string{
		`stdlog_test.go:20: DEBUG: isn't this "function"...`,
		`stdlog_test.go:21: DEBUG: ...fun?`,
	}

	for i, line := range lines {
		prefix := expectedVariablePrefix.FindString(line)
		assert.NotEmpty(t, prefix)
		line = line[len(prefix):]
		assert.Equal(t, expectedLinesWithoutPrefix[i], line)
	}
}
