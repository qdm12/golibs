package logging

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Logger,ParentLogger

type Logger interface {
	// Debug logs with the Debug level.
	Debug(s string)
	// Info logs with the Info level.
	Info(s string)
	// Warnf logs with the Warning level.
	Warn(s string)
	// Error logs with the Error level.
	Error(s string)
}

type ParentLogger interface {
	Logger
	// NewChild creates a child logger with the same writer as
	// the current logger and with the settings given.
	// It should be used to have thread safety on the same writer.
	// Note that the Writer setting is ignored.
	NewChild(settings Settings) ParentLogger
}
