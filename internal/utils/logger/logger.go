package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.SugaredLogger

func Init(env string, level string, logFilePath string) {
	var cfg zap.Config
	var err error

	switch env {
	case "production":
		cfg = zap.NewProductionConfig()
		cfg.OutputPaths = []string{logFilePath, "stderr"}
		cfg.ErrorOutputPaths = []string{logFilePath, "stderr"}
	case "development":
		cfg = zap.NewDevelopmentConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		log.Fatalf("Unknown environment: %s", level)
	}

	if err := cfg.Level.UnmarshalText([]byte(level)); err != nil {
		log.Fatalf("Invalid log level: %s", level)
	}

	zapLogger, err := cfg.Build()
	if err != nil {
		log.Fatalf("Cannot initialize zap logger: %v", err)
	}

	Log = zapLogger.Sugar()

	Log.Infow("Logger initialized",
		"env", env,
		"level", level,
		"logFilePath", logFilePath,
	)
}

func Errorw(msg string, keysAndValues ...any) {
	if Log != nil {
		Log.Errorw(msg, keysAndValues...)
	}
}

func Infow(msg string, keysAndValues ...any) {
	if Log != nil {
		Log.Infow(msg, keysAndValues...)
	}
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
