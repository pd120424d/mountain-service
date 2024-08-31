package utils

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
	Sync() error
}
type zapLogger struct {
	logger *zap.Logger
}

func (z *zapLogger) Debug(msg string, fields ...zap.Field) {
	z.logger.Debug(msg, fields...)
}

func (z *zapLogger) Info(msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

func (z *zapLogger) Warn(msg string, fields ...zap.Field) {
	z.logger.Warn(msg, fields...)
}

func (z *zapLogger) Error(msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

func (z *zapLogger) Fatal(msg string, fields ...zap.Field) {
	z.logger.Fatal(msg, fields...)
}

func (z *zapLogger) Debugf(format string, args ...interface{}) {
	z.logger.Debug(fmt.Sprintf(format, args...))
}

func (z *zapLogger) Infof(format string, args ...interface{}) {
	z.logger.Info(fmt.Sprintf(format, args...))
}

func (z *zapLogger) Warnf(format string, args ...interface{}) {
	z.logger.Warn(fmt.Sprintf(format, args...))
}

func (z *zapLogger) Errorf(format string, args ...interface{}) {
	z.logger.Error(fmt.Sprintf(format, args...))
}

func (z *zapLogger) Fatalf(format string, args ...interface{}) {
	z.logger.Fatal(fmt.Sprintf(format, args...))
}

func (z *zapLogger) Sync() error {
	return z.logger.Sync()
}

func NewLogger() Logger {
	var zapConfig zap.Config

	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		zapConfig = zap.NewDevelopmentConfig()
	case "PROD":
		zapConfig = zap.NewProductionConfig()
	default:
		zapConfig = zap.NewDevelopmentConfig()
	}

	logger, err := zapConfig.Build()
	if err != nil {
		panic(err)
	}

	return &zapLogger{logger}
}
