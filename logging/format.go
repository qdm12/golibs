package logging

import "fmt"

// Debugf formats and logs with the Debug level
func Debugf(format string, a ...interface{}) {
	Debug(fmt.Sprintf(format, a...))
}

// Infof formats and logs with the Info level
func Infof(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a...))
}

// Warnf formats and logs with the Warning level
func Warnf(format string, a ...interface{}) {
	Warn(fmt.Sprintf(format, a...))
}

// Errorf formats and logs with the Error level
func Errorf(format string, a ...interface{}) {
	Error(fmt.Sprintf(format, a...))
}
