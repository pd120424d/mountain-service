package utils

import (
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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
	RequestLogger() gin.HandlerFunc
}

type zapLogger struct {
	mutex       sync.Mutex
	logger      *zap.Logger
	file        *os.File
	svcName     string
	currentDate string
}

func NewLogger(svcName string) (Logger, error) {
	z := &zapLogger{svcName: svcName}
	if err := z.rotate(); err != nil {
		return nil, err
	}
	return z, nil
}

func (z *zapLogger) rotate() error {
	z.mutex.Lock()
	defer z.mutex.Unlock()

	currentDate := time.Now().Format("2006-01-02")
	if z.currentDate == currentDate {
		return nil // No rotation needed
	}

	// Close old file if exists
	if z.file != nil {
		_ = z.file.Close()
	}

	// New log file
	filename := fmt.Sprintf("/var/log/%s.%s.log", z.svcName, currentDate)
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	logger := newLoggerFromWriter(file)

	if z.logger != nil {
		_ = z.logger.Sync()
	}

	z.file = file
	z.logger = logger
	z.currentDate = currentDate

	return nil
}

// Logging methods with rotation check
func (z *zapLogger) Debug(msg string, fields ...zap.Field) {
	_ = z.rotate()
	z.logger.Debug(msg, fields...)
}

func (z *zapLogger) Info(msg string, fields ...zap.Field) {
	_ = z.rotate()
	z.logger.Info(msg, fields...)
}

func (z *zapLogger) Warn(msg string, fields ...zap.Field) {
	_ = z.rotate()
	z.logger.Warn(msg, fields...)
}

func (z *zapLogger) Error(msg string, fields ...zap.Field) {
	_ = z.rotate()
	z.logger.Error(msg, fields...)
}

func (z *zapLogger) Fatal(msg string, fields ...zap.Field) {
	_ = z.rotate()
	z.logger.Fatal(msg, fields...)
}

func (z *zapLogger) Debugf(format string, args ...interface{}) {
	_ = z.rotate()
	z.logger.Debug(fmt.Sprintf(format, args...))
}

func (z *zapLogger) Infof(format string, args ...interface{}) {
	_ = z.rotate()
	z.logger.Info(fmt.Sprintf(format, args...))
}

func (z *zapLogger) Warnf(format string, args ...interface{}) {
	_ = z.rotate()
	z.logger.Warn(fmt.Sprintf(format, args...))
}

func (z *zapLogger) Errorf(format string, args ...interface{}) {
	_ = z.rotate()
	z.logger.Error(fmt.Sprintf(format, args...))
}

func (z *zapLogger) Fatalf(format string, args ...interface{}) {
	_ = z.rotate()
	z.logger.Fatal(fmt.Sprintf(format, args...))
}

func (z *zapLogger) Sync() error {
	z.mutex.Lock()
	defer z.mutex.Unlock()
	if z.logger != nil {
		return z.logger.Sync()
	}
	return nil
}

func (z *zapLogger) WithName(name string) Logger {
	_ = z.rotate()
	named := z.logger.Named(name)
	return &zapLogger{
		logger:      named,
		file:        z.file,
		svcName:     z.svcName,
		currentDate: z.currentDate,
	}
}

func NewNamedLogger(baseLogger Logger, name string) Logger {
	return baseLogger.WithName(name)
}

func NewTestLogger() Logger {
	return &zapLogger{logger: newLoggerFromWriter(os.Stdout)}
}

func newLoggerFromWriter(writer io.Writer) *zap.Logger {
	writeSyncer := zapcore.AddSync(writer)

	var encoderConfig zapcore.EncoderConfig
	level := os.Getenv("LOG_LEVEL")
	if level == "DEBUG" {
		encoderConfig = zap.NewProductionEncoderConfig()
		encoderConfig.LevelKey = "level"
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"
	encoderConfig.NameKey = "logger"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		writeSyncer,
		zapcore.DebugLevel,
	)

	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

func (l *zapLogger) RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// Process request
		c.Next()

		// After request
		duration := time.Since(start)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		l.Infof("[%d] %s %s from %s (%s)", status, method, path, clientIP, duration)
	}
}
