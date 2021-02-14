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
	// NewChild creates a child logger with the options given
	// and is thread safe to use with other child loggers.
	// It should really only be used to create another logger
	// to the same io.Writer such that it is thread safe.
	NewChild(options ...Option) Logger
}

// New creates a new logger.
// It should only be called once per writer (options.Writer), child loggers
// can be created using the WithOptions method.
func New(loggerType Type, options ...Option) Logger {
	settings := newSettings(options...)

	switch loggerType {
	case StdLog:
		return newStdLog(settings)
	default:
		panic("logger type " + loggerType + " not supported")
	}
}
