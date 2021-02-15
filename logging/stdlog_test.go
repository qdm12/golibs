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

	postProcess := func(line string) string {
		return line + " (postprocessed)"
	}

	logger := New(StdLog, SetWriter(buffer), SetCaller(CallerShort), SetPostProcess(postProcess))

	logger.Debug("isn't this %q...", "function")
	logger.Debug("...fun?")

	result := buffer.String()
	buffer.Reset()

	result = strings.TrimSuffix(result, "\n")

	lines := strings.Split(result, "\n")
	require.Len(t, lines, 2)

	expectedVariablePrefix := regexp.MustCompile(`2[0-9]{3}/[0-1][0-9]/[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] `)

	expectedLinesWithoutPrefix := []string{
		`stdlog_test.go:24: DEBUG: isn't this "function"... (postprocessed)`,
		`stdlog_test.go:25: DEBUG: ...fun? (postprocessed)`,
	}

	for i, line := range lines {
		prefix := expectedVariablePrefix.FindString(line)
		assert.NotEmpty(t, prefix)
		line = line[len(prefix):]
		assert.Equal(t, expectedLinesWithoutPrefix[i], line)
	}
}
