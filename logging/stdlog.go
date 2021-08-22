package logging

import (
	"log"
	"sync"
)

var _ ParentLogger = (*StdLogger)(nil)
var _ Logger = (*StdLogger)(nil)

type StdLogger struct {
	logImpl     *log.Logger
	settings    Settings
	writerMutex *sync.Mutex
}

// New creates a new logger using the standard library logger.
// It should only be called once at most per writer (settings.Writer).
// If you want to create more loggers with different settings for the
// same writer, child loggers can be created using the NewChild method,
// to ensure thread safety on the same writer.
func New(settings Settings) *StdLogger {
	flags := log.Ldate | log.Ltime
	if settings.Caller == CallerShort {
		flags |= log.Lshortfile
	}
	settings.setDefaults()

	logImpl := log.New(settings.Writer, "", flags)

	return &StdLogger{
		logImpl:     logImpl,
		settings:    settings,
		writerMutex: &sync.Mutex{},
	}
}

func (l *StdLogger) NewChild(settings Settings) ParentLogger {
	settings.setEmptyValuesWith(l.settings)
	settings.setDefaults()

	flags := log.Ldate | log.Ltime
	if settings.Caller == CallerShort {
		flags |= log.Lshortfile
	}

	logImpl := log.New(l.logImpl.Writer(), "", flags)

	return &StdLogger{
		logImpl:     logImpl,
		settings:    settings,
		writerMutex: l.writerMutex,
	}
}

func (l *StdLogger) log(logLevel Level, s string) {
	if l.settings.Level > logLevel {
		return
	}

	// Line is computed here to avoid blocking child loggers with its
	// computing and the eventual preprocess function.
	line := logLevel.String() + " " + formatWithSettings(l.settings, s)

	l.writerMutex.Lock()
	defer l.writerMutex.Unlock()

	const callDepth = 3
	_ = l.logImpl.Output(callDepth, line)
}

func (l *StdLogger) Debug(s string) { l.log(LevelDebug, s) }
func (l *StdLogger) Info(s string)  { l.log(LevelInfo, s) }
func (l *StdLogger) Warn(s string)  { l.log(LevelWarn, s) }
func (l *StdLogger) Error(s string) { l.log(LevelError, s) }
