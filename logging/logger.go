package logging

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logging interface {
	// Sync synchronizes the buffer to ensure everything is printed out
	Sync() error
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
}

type logging struct {
	logger *zap.Logger
}

func NewLogging(encoding Encoding, level Level, nodeID int) (l Logging, err error) {
	var zapLevel zapcore.Level
	zapLevel.UnmarshalText([]byte(string(level)))
	logger, err := zap.Config{
		Level:    zap.NewAtomicLevelAt(zapLevel),
		Encoding: string(encoding),
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
	}.Build()
	if err != nil {
		return nil, fmt.Errorf("logger initialization failed: %w", err)
	}
	if nodeID != -1 {
		logger = logger.With(zap.Int("node_id", nodeID))
	}
	return &logging{logger}, nil
}

func (l *logging) Sync() error {
	return l.logger.Sync()
}

// NewEmptyLogging returns a logging that does not print anything
func NewEmptyLogging() (l Logging, err error) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.OutputPaths, loggerConfig.ErrorOutputPaths = nil, nil
	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}
	return &logging{logger}, nil
}
