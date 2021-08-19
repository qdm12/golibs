package logging

import (
	"log"
	"sync"
)

type stdLogger struct {
	logImpl  *log.Logger
	settings Settings
	// isConcurrent and writerMutex are only used
	// when more loggers are created using the
	// NewChild method to avoid writing to the
	// same writer at the same time.
	isConcurrent bool
	writerMutex  *sync.Mutex
}

// New creates a new logger using the standard library logger.
// It should only be called once at most per writer (settings.Writer).
func New(settings Settings) Logger {
	return NewParent(settings)
}

// NewParent creates a new logger using the standard library logger.
// It should only be called once at most per writer (settings.Writer).
// If you want to create more loggers with different settings for the
// same writer, child loggers can be created using the NewChild method,
// to ensure thread safety on the same writer.
func NewParent(settings Settings) ParentLogger {
	flags := log.Ldate | log.Ltime
	if settings.Caller == CallerShort {
		flags |= log.Lshortfile
	}
	settings.setDefaults()

	logImpl := log.New(settings.Writer, "", flags)

	return &stdLogger{
		logImpl:     logImpl,
		settings:    settings,
		writerMutex: &sync.Mutex{},
	}
}

func (l *stdLogger) NewChild(settings Settings) ParentLogger {
	newSettings := l.settings
	newSettings.setEmptyValuesWith(settings)
	newSettings.setDefaults()

	l.isConcurrent = true

	flags := log.Ldate | log.Ltime
	if newSettings.Caller == CallerShort {
		flags |= log.Lshortfile
	}

	logImpl := log.New(l.logImpl.Writer(), "", flags)

	return &stdLogger{
		logImpl:      logImpl,
		settings:     newSettings,
		isConcurrent: true,
		writerMutex:  l.writerMutex,
	}
}

func (l *stdLogger) log(logLevel Level, s string) {
	if l.settings.Level > logLevel {
		return
	}

	// Line is computed here to avoid blocking child loggers with its
	// computing and the eventual preprocess function.
	line := logLevel.String() + " " + formatWithSettings(l.settings, s)

	if l.isConcurrent {
		l.writerMutex.Lock()
		defer l.writerMutex.Unlock()
	}

	const callDepth = 3
	_ = l.logImpl.Output(callDepth, line)
}

func (l *stdLogger) Debug(s string) { l.log(LevelDebug, s) }
func (l *stdLogger) Info(s string)  { l.log(LevelInfo, s) }
func (l *stdLogger) Warn(s string)  { l.log(LevelWarn, s) }
func (l *stdLogger) Error(s string) { l.log(LevelError, s) }
