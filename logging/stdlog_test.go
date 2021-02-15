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

	preprocess := func(line string) string {
		return line + " (preprocessed)"
	}

	logger := New(StdLog, SetLevel(LevelDebug), SetWriter(buffer),
		SetCaller(CallerShort), SetPreProcess(preprocess), SetPrefix("server: "))

	logger.Debug("isn't this %q...", "function")
	logger.Debug("...fun?")

	result := buffer.String()
	buffer.Reset()

	result = strings.TrimSuffix(result, "\n")

	lines := strings.Split(result, "\n")
	require.Len(t, lines, 2)

	expectedVariablePrefix := regexp.MustCompile(`2[0-9]{3}/[0-1][0-9]/[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] `)

	expectedLinesWithoutPrefix := []string{
		`stdlog_test.go:25: DEBUG server: isn't this "function"... (preprocessed)`,
		`stdlog_test.go:26: DEBUG server: ...fun? (preprocessed)`,
	}

	for i, line := range lines {
		prefix := expectedVariablePrefix.FindString(line)
		assert.NotEmpty(t, prefix)
		line = line[len(prefix):]
		assert.Equal(t, expectedLinesWithoutPrefix[i], line)
	}
}
