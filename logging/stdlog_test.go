package logging

import (
	"bytes"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_StdLog_Debug(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBuffer(nil)

	preprocess := func(line string) string {
		return line + " (preprocessed)"
	}

	logger := New(Settings{
		Level:      LevelDebug,
		Writer:     buffer,
		Caller:     CallerShort,
		PreProcess: preprocess,
		Prefix:     "server: ",
	})

	logger.Debug("isn't this \"function\"...")
	logger.Debug("...fun?")

	result := buffer.String()
	buffer.Reset()

	result = strings.TrimSuffix(result, "\n")

	lines := strings.Split(result, "\n")
	require.Len(t, lines, 2)

	expectedVariablePrefix := regexp.MustCompile(`2[0-9]{3}/[0-1][0-9]/[0-3][0-9] [0-2][0-9]:[0-5][0-9]:[0-5][0-9] `)

	expectedLinesWithoutPrefix := []string{
		`stdlog_test.go:30: DEBUG server: isn't this "function"... (preprocessed)`,
		`stdlog_test.go:31: DEBUG server: ...fun? (preprocessed)`,
	}

	for i, line := range lines {
		prefix := expectedVariablePrefix.FindString(line)
		assert.NotEmpty(t, prefix)
		line = line[len(prefix):]
		assert.Equal(t, expectedLinesWithoutPrefix[i], line)
	}
}

func Test_StdLogger_PatchLevel(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBuffer(nil)

	logger := New(Settings{
		Level:  LevelError,
		Writer: buffer,
	})

	logger.PatchLevel(LevelInfo)

	expectedSettings := Settings{
		Level:  LevelInfo,
		Writer: buffer,
	}
	assert.Equal(t, expectedSettings, logger.settings)
}

func Test_StdLogger_PatchPrefix(t *testing.T) {
	t.Parallel()

	buffer := bytes.NewBuffer(nil)

	logger := New(Settings{
		Prefix: "prefix",
		Writer: buffer,
	})

	logger.PatchPrefix("new prefix")

	expectedSettings := Settings{
		Prefix: "new prefix",
		Writer: buffer,
	}
	assert.Equal(t, expectedSettings, logger.settings)
}
