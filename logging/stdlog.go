package logging

import (
	"log"
	"sync"
)

type stdLogger struct {
	logImpl  *log.Logger
	settings settings
	// hasChild and writerMutex are only used
	// when more loggers are created using the
	// NewChild method to avoid writing to the
	// same writer at the same time.
	hasChild    bool
	writerMutex sync.Mutex
}

func newStdLog(settings settings) Logger {
	flags := log.Ldate | log.Ltime
	if settings.caller == CallerShort {
		flags |= log.Lshortfile
	}

	logImpl := log.New(settings.writer, settings.prefix, flags)

	return &stdLogger{
		logImpl: logImpl,
	}
}

func (l *stdLogger) NewChild(options ...Option) Logger {
	settings := newSettings(options...)

	l.hasChild = true

	flags := log.Ldate | log.Ltime
	if settings.caller == CallerShort {
		flags |= log.Lshortfile
	}

	logImpl := log.New(settings.writer, settings.prefix, flags)

	return &stdLogger{
		logImpl: logImpl,
	}
}

func (l *stdLogger) log(logLevel Level, args ...interface{}) {
	if l.settings.level < logLevel {
		return
	}

	if l.hasChild {
		l.writerMutex.Lock()
		defer l.writerMutex.Unlock()
	}

	const callDepth = 3
	l.logImpl.Output(callDepth, logLevel.String()+": "+
		formatWithSettings(l.settings, args...))
}

func (l *stdLogger) Debug(args ...interface{}) { l.log(LevelDebug, args...) }
func (l *stdLogger) Info(args ...interface{})  { l.log(LevelInfo, args...) }
func (l *stdLogger) Warn(args ...interface{})  { l.log(LevelWarn, args...) }
func (l *stdLogger) Error(args ...interface{}) { l.log(LevelError, args...) }
