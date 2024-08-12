package utils

import (
	"go.uber.org/zap"
)

type Logger struct {
	SugaredLogger *zap.SugaredLogger
}

var (
	prodLogger    *zap.SugaredLogger
	stagingLogger *zap.SugaredLogger
)

func init() {
	prodConfig := zap.NewProductionConfig()
	prodLogger = newLogger(prodConfig)

	stagingConfig := zap.NewDevelopmentConfig()
	stagingLogger = newLogger(stagingConfig)
}

func newLogger(cfg zap.Config) *zap.SugaredLogger {
	logger, err := cfg.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}

// GetLogger is a constructor for the Logger that Fx can use
func GetLogger(env string) func() *Logger {
	return func() *Logger {
		var logger *zap.SugaredLogger
		switch env {
		case "production":
			logger = prodLogger
		case "staging":
			logger = stagingLogger
		default:
			logger = stagingLogger
		}
		return &Logger{SugaredLogger: logger}
	}
}
