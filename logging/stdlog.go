package logging

import (
	"log"
	"sync"
)

var _ ParentLogger = (*StdLogger)(nil)
var _ Logger = (*StdLogger)(nil)

type StdLogger struct {
	logImpl  *log.Logger
	settings Settings
	mutex    *sync.Mutex
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
		logImpl:  logImpl,
		settings: settings,
		mutex:    &sync.Mutex{},
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
		logImpl:  logImpl,
		settings: settings,
		mutex:    l.mutex,
	}
}

func (l *StdLogger) log(logLevel Level, s string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.settings.Level > logLevel {
		return
	}

	// Line is computed here to avoid blocking child loggers with its
	// computing and the eventual preprocess function.
	line := logLevel.String() + " " + formatWithSettings(l.settings, s)

	const callDepth = 3
	_ = l.logImpl.Output(callDepth, line)
}

func (l *StdLogger) Debug(s string) { l.log(LevelDebug, s) }
func (l *StdLogger) Info(s string)  { l.log(LevelInfo, s) }
func (l *StdLogger) Warn(s string)  { l.log(LevelWarn, s) }
func (l *StdLogger) Error(s string) { l.log(LevelError, s) }

// PatchLevel changes the level of the logger.
// Note it does not change the level of child loggers.
func (l *StdLogger) PatchLevel(level Level) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.settings.Level = level
}

// PatchPrefix changes the prefix of the logger.
// Note it does not change the prefix of child loggers.
func (l *StdLogger) PatchPrefix(prefix string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.settings.Prefix = prefix
}
