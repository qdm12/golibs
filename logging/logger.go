package logging

//go:generate mockgen -destination=mock_$GOPACKAGE/$GOFILE . Logger

type Logger interface {
	// Debug formats and logs with the Debug level.
	Debug(args ...interface{})
	// Info formats and logs with the Info level.
	Info(args ...interface{})
	// Warnf formats and logs with the Warning level.
	Warn(args ...interface{})
	// Error formats and logs with the Error level.
	Error(args ...interface{})
	// NewChild creates a child logger with the same writer as
	// the current logger and with the settings given.
	// It should be used to have thread safety on the same writer.
	// Note that the Writer setting is ignored.
	NewChild(settings Settings) Logger
}

// New creates a new logger.
// It should only be called once per writer (options.Writer), child loggers
// can be created using the NewChild method.
func New(loggerType Type, settings Settings) Logger {
	settings.setDefaults()

	switch loggerType {
	case StdLog:
		return newStdLog(settings)
	default:
		panic("logger type " + loggerType + " not supported")
	}
}
