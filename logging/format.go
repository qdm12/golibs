package logging

import (
	"github.com/qdm12/golibs/format"
	"go.uber.org/zap/zapcore"
)

// Debug formats and logs with the Debug level.
func (l *logger) Debug(args ...interface{}) {
	l.log(l.zapLogger.Debug, args...)
}

// Info formats and logs with the Info level.
func (l *logger) Info(args ...interface{}) {
	l.log(l.zapLogger.Info, args...)
}

// Warnf formats and logs with the Warning level.
func (l *logger) Warn(args ...interface{}) {
	l.log(l.zapLogger.Warn, args...)
}

// Error formats and logs with the Error level.
func (l *logger) Error(args ...interface{}) {
	l.log(l.zapLogger.Error, args...)
}

type zapLogFn func(s string, fields ...zapcore.Field)

func (l *logger) log(logFn zapLogFn, args ...interface{}) {
	s := l.prefix + format.ArgsToString(args...)
	if l.color != nil {
		s = l.color.Sprintf(s)
	}
	logFn(s)
}
