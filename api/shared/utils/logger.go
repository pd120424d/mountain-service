package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/logging"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/api/option"
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

	// GCP client for Cloud Logging flush
	gcpClient *logging.Client
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
		z.file = nil
	}

	isKubernetes := os.Getenv("KUBERNETES_SERVICE_HOST") != ""
	forceStdoutOnly := os.Getenv("LOG_TO_STDOUT_ONLY") == "true" || isKubernetes
	enableFile := os.Getenv("LOG_TO_FILE") == "true" && !forceStdoutOnly

	var writers []io.Writer

	// Setup file writer (local/dev) if enabled
	if enableFile {
		logDir := os.Getenv("LOG_DIR")
		if logDir == "" {
			logDir = "/var/log"
		}
		// Create log directory if it doesn't exist
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}
		// New log file (date-based)
		filename := fmt.Sprintf("%s/%s.%s.log", logDir, z.svcName, currentDate)
		file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		z.file = file
		writers = append(writers, file)
	}

	if forceStdoutOnly || os.Getenv("LOG_TO_STDOUT") == "true" || isKubernetes {
		writers = append(writers, os.Stdout)
	}
	if len(writers) == 0 {
		writers = append(writers, os.Stdout)
	}

	// Build base zap core (JSON) with combined writer(s)
	encoderConfig := buildEncoderConfig()
	baseCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.AddSync(io.MultiWriter(writers...)),
		zapcore.DebugLevel,
	)

	cores := []zapcore.Core{baseCore}

	// Optionally add Google Cloud Logging cores
	if os.Getenv("GCP_LOGGING_ENABLED") == "true" {
		projectID := firstNonEmpty(os.Getenv("GCP_PROJECT_ID"), os.Getenv("FIREBASE_PROJECT_ID"))
		logID := os.Getenv("GCP_LOG_ID")
		if logID == "" {
			logID = z.svcName
		}
		if projectID != "" {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			// Prefer separate credentials for logging if provided
			var client *logging.Client
			var err error
			if credPath := os.Getenv("GCP_LOGGING_CREDENTIALS_PATH"); credPath != "" {
				client, err = logging.NewClient(ctx, projectID, option.WithCredentialsFile(credPath))
			} else {
				client, err = logging.NewClient(ctx, projectID)
			}
			if err == nil {
				z.gcpClient = client
				gcpLogger := client.Logger(logID)

				// One writer per severity to preserve levels in Cloud Logging
				wDebug := gcpLogger.StandardLogger(logging.Debug).Writer()
				wInfo := gcpLogger.StandardLogger(logging.Info).Writer()
				wWarn := gcpLogger.StandardLogger(logging.Warning).Writer()
				wError := gcpLogger.StandardLogger(logging.Error).Writer()

				enc := zapcore.NewJSONEncoder(encoderConfig)
				cores = append(cores,
					zapcore.NewCore(enc, zapcore.AddSync(wDebug), exactLevelEnabler{level: zapcore.DebugLevel}),
					zapcore.NewCore(enc, zapcore.AddSync(wInfo), exactLevelEnabler{level: zapcore.InfoLevel}),
					zapcore.NewCore(enc, zapcore.AddSync(wWarn), exactLevelEnabler{level: zapcore.WarnLevel}),
					zapcore.NewCore(enc, zapcore.AddSync(wError), errorAndAboveEnabler{}),
				)
			}
		}
	}

	logger := zap.New(zapcore.NewTee(cores...), zap.AddCaller(), zap.AddCallerSkip(1))

	if z.logger != nil {
		_ = z.logger.Sync()
	}

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
	// Flush zap buffers
	if z.logger != nil {
		_ = z.logger.Sync()
	}
	// Flush Cloud Logging, if configured
	if z.gcpClient != nil {
		return z.gcpClient.Close()
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
		gcpClient:   z.gcpClient,
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
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(buildEncoderConfig()),
		writeSyncer,
		zapcore.DebugLevel,
	)
	return zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
}

// RequestLogger returns a Gin middleware that logs request info
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
		requestID := RequestIDFromContext(c.Request.Context())

		if requestID != "" {
			l.Info("http_request",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
				zap.Duration("duration", duration),
				zap.String("request_id", requestID),
			)
		} else {
			l.Info("http_request",
				zap.Int("status", status),
				zap.String("method", method),
				zap.String("path", path),
				zap.String("client_ip", clientIP),
				zap.Duration("duration", duration),
			)
		}
	}
}

func buildEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"
	encoderConfig.NameKey = "logger"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	return encoderConfig
}

type exactLevelEnabler struct{ level zapcore.Level }

func (e exactLevelEnabler) Enabled(l zapcore.Level) bool { return l == e.level }

type errorAndAboveEnabler struct{}

func (errorAndAboveEnabler) Enabled(l zapcore.Level) bool { return l >= zapcore.ErrorLevel }

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
