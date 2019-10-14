package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// AssertErrosEqual asserts errors match, so that they are both nil or their messages
// are both equal
func AssertErrorsEqual(t *testing.T, expected, actual error) (success bool) {
	if expected == nil {
		return assert.Nil(t, actual)
	}
	return assert.EqualError(t, actual, expected.Error())
}

// SetDefaultLoggerToEmpty sets Zap's default logger to output to no stream
func SetDefaultLoggerToEmpty() (restore func()) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths, loggerConfig.ErrorOutputPaths = nil, nil
	logger, err := loggerConfig.Build()
	if err != nil {
		panic(err)
	}
	return zap.ReplaceGlobals(logger)
}
