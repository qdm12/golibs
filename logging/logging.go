package logging

import (
	"github.com/qdm12/golibs/params"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Log(message string) {
	zap.L().Info(message)
}

func Warn(message string) {
	zap.L().Warn(message)
}

func Error(message string, err error) {
	zap.L().Error(message, zap.Error(err))
}

func Sync() error {
	return zap.L().Sync()
}

func InitLogger() {
	encoding, err := params.GetLoggerEncoding()
	if err != nil {
		Error("logger initialization failed", err)
	}
	level, err := params.GetLoggerLevel()
	if err != nil {
		Error("logger initialization failed", err)
	}
	nodeID, err := params.GetNodeID()
	if err != nil {
		Error("logger initialization failed", err)
	}
	config := zap.Config{
		Level:    zap.NewAtomicLevelAt(level),
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
		Error("logger initialization failed", err)
	}
	logger = logger.With(zap.Int("node_id", nodeID))
	zap.ReplaceGlobals(logger)
}
