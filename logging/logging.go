package logging

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Level is the level of the logger
type Level zapcore.Level

const (
	InfoLevel = iota
	WarnLevel
	ErrorLevel
)

// InitLogger sets up the global logger using the parameters given.
// If nodeID is set to -1, it is not added to the log messages.
func InitLogger(encoding string, level Level, nodeID int) {
	config := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapcore.Level(level)),
		Encoding: encoding,
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			MessageKey:     "msg",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
	}
	logger, err := config.Build()
	if err != nil {
		Errorf("logger initialization failed: %s", err)
	}
	if nodeID != -1 {
		logger = logger.With(zap.Int("node_id", nodeID))
	}
	zap.ReplaceGlobals(logger)
}

// Sync synchronizes the buffer to ensure everything is printed out
func Sync() error {
	return zap.L().Sync()
}

// SetLoggerToEmpty sets the logger stdout and stderr to nil for testing
func SetLoggerToEmpty() (restore func(), err error) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths, loggerConfig.ErrorOutputPaths = nil, nil
	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}
	return zap.ReplaceGlobals(logger), nil
}
