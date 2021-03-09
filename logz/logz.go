package logz

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogConfig() *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.MessageKey = "message"
	encoderConfig.CallerKey = "caller"

	config := zap.NewProductionConfig()
	var logLevel zapcore.Level
	switch viper.GetString("log.level") {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.InfoLevel
	}
	config.Level = zap.NewAtomicLevelAt(logLevel)
	if viper.GetString("log.env") == "dev" {
		config.Encoding = "console"
	} else {
		config.Encoding = "json"
	}
	config.EncoderConfig = encoderConfig

	logger, err := config.Build()
	if err != nil {
		log.Fatal("fatal error config logger")
	}
	return logger
}

func ExecutionTime(start time.Time, name string, l *zap.Logger) {
	elapse := time.Since(start)
	l.With(zap.Duration("duration", elapse)).Info(fmt.Sprintf("%s took %s", name, elapse))
}
