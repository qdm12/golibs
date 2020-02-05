package logging

import "github.com/qdm12/golibs/format"

// Debug formats and logs with the Debug level
func (l *logger) Debug(args ...interface{}) {
	l.zapLogger.Debug(l.prefix + format.ArgsToString(args...))
}

// Info formats and logs with the Info level
func (l *logger) Info(args ...interface{}) {
	l.zapLogger.Info(l.prefix + format.ArgsToString(args...))
}

// Warnf formats and logs with the Warning level
func (l *logger) Warn(args ...interface{}) {
	l.zapLogger.Warn(l.prefix + format.ArgsToString(args...))
}

// Error formats and logs with the Error level
func (l *logger) Error(args ...interface{}) {
	l.zapLogger.Error(l.prefix + format.ArgsToString(args...))
}
