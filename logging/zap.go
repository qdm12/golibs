package logging

import "go.uber.org/zap/zapcore"

//go:generate mockgen -destination=mockZap_test.go -package=logging github.com/qdm12/golibs/logging Zap
type Zap interface {
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
	Sync() error
}
