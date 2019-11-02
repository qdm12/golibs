package logging

import "go.uber.org/zap"

// Debug logs a string with the Debug level
func Debug(message string) {
	zap.L().Debug(message)
}

// Info logs a string with the Info level
func Info(message string) {
	zap.L().Info(message)
}

// Warn logs a string with the Warning level
func Warn(message string) {
	zap.L().Warn(message)
}

// Error logs a string with the Error level
func Error(message string) {
	zap.L().Error(message)
}

// Werr logs an error with the Warning level
func Werr(err error) {
	zap.L().Warn(err.Error())
}

// Err logs an error with the Error level
func Err(err error) {
	zap.L().Error(err.Error())
}
