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
// It should only be called once per writer (options.Writer), child loggers
// can be created using the NewChild method.
func New(settings Settings) Logger {
	flags := log.Ldate | log.Ltime
	if settings.Caller == CallerShort {
		flags |= log.Lshortfile
	}

	logImpl := log.New(settings.Writer, "", flags)

	return &stdLogger{
		logImpl:     logImpl,
		settings:    settings,
		writerMutex: &sync.Mutex{},
	}
}

func (l *stdLogger) NewChild(settings Settings) Logger {
	settings.setDefaults()

	l.isConcurrent = true

	flags := log.Ldate | log.Ltime
	if settings.Caller == CallerShort {
		flags |= log.Lshortfile
	}

	logImpl := log.New(l.logImpl.Writer(), "", flags)

	return &stdLogger{
		logImpl:      logImpl,
		settings:     settings,
		isConcurrent: true,
		writerMutex:  l.writerMutex,
	}
}

func (l *stdLogger) log(logLevel Level, args ...interface{}) {
	if l.settings.Level > logLevel {
		return
	}

	// Line is computed here to avoid blocking child loggers with its
	// computing and the eventual preprocess function.
	line := logLevel.String() + " " + formatWithSettings(l.settings, args...)

	if l.isConcurrent {
		l.writerMutex.Lock()
		defer l.writerMutex.Unlock()
	}

	const callDepth = 3
	_ = l.logImpl.Output(callDepth, line)
}

func (l *stdLogger) Debug(args ...interface{}) { l.log(LevelDebug, args...) }
func (l *stdLogger) Info(args ...interface{})  { l.log(LevelInfo, args...) }
func (l *stdLogger) Warn(args ...interface{})  { l.log(LevelWarn, args...) }
func (l *stdLogger) Error(args ...interface{}) { l.log(LevelError, args...) }
