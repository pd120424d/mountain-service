package utils

import (
	"fmt"
	"io"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	WithName(name string) Logger
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

func (z *zapLogger) WithName(name string) Logger {
	namedLogger := z.logger.Named(name)
	return &zapLogger{logger: namedLogger}
}

func NewLogger(svcName string) (Logger, error) {
	logFile := fmt.Sprintf("/var/log/%s.%s.log", svcName, time.Now().Format("2006-01-02"))

	// Open the file in append mode, create it if it doesn't exist
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return newLoggerFromFile(file), nil
}

func NewNamedLogger(baseLogger Logger, name string) Logger {
	return baseLogger.WithName(name)
}

func NewTestLogger() Logger {
	// Use stdout for testing
	return newLoggerFromFile(os.Stdout)
}

func newLoggerFromFile(file io.Writer) Logger {
	writeSyncer := zapcore.AddSync(file)

	var encoderConfig zapcore.EncoderConfig

	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	case "PROD":
		encoderConfig = zap.NewProductionEncoderConfig()
	default:
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	logLevel := zapcore.DebugLevel

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create the core logger
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // Use JSON format
		writeSyncer,
		logLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return &zapLogger{logger: logger}
}
