package logging

import (
	"github.com/qdm12/golibs/params"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger sets up the global logger using the environment variables
// LOGENCODING, LOGLEVEL and NODEID
func InitLogger() {
	encoding, err := params.GetLoggerEncoding()
	if err != nil {
		Errorf("logger initialization failed: %s", err)
	}
	level, err := params.GetLoggerLevel()
	if err != nil {
		Errorf("logger initialization failed: %s", err)
	}
	nodeID, err := params.GetNodeID()
	if err != nil {
		Errorf("logger initialization failed: %s", err)
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
		Errorf("logger initialization failed: %s", err)
	}
	logger = logger.With(zap.Int("node_id", nodeID))
	zap.ReplaceGlobals(logger)
}

// Sync synchronizes the buffer to ensure everything is printed out
func Sync() error {
	return zap.L().Sync()
}
