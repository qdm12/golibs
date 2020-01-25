package logging

import "github.com/qdm12/golibs/format"

// Debug formats and logs with the Debug level
func (l *logging) Debug(args ...interface{}) {
	l.logger.Debug(format.ArgsToString(args))
}

// Info formats and logs with the Info level
func (l *logging) Info(args ...interface{}) {
	l.logger.Info(format.ArgsToString(args))
}

// Warnf formats and logs with the Warning level
func (l *logging) Warn(args ...interface{}) {
	l.logger.Warn(format.ArgsToString(args))
}

// Error formats and logs with the Error level
func (l *logging) Error(args ...interface{}) {
	l.logger.Error(format.ArgsToString(args))
}
