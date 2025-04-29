package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	debugLevel = "debug"
	infoLevel  = "info"
	warnLevel  = "warn"
	errorLevel = "fatal"
	panicLevel = "panic"
)

func NewClientZapLogger(logLevel string, clientID string) *zap.Logger {
	return NewZapLogger(logLevel).With(zap.String("client-id", clientID))
}

func NewZapLogger(logLevel string) *zap.Logger {
	var level zap.AtomicLevel

	switch logLevel {
	case debugLevel:
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case infoLevel:
		level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case warnLevel:
		level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case errorLevel:
		level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	case panicLevel:
		level = zap.NewAtomicLevelAt(zap.PanicLevel)
	}

	encoderConfig := zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "trace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	config := zap.Config{
		Level:            level,
		Development:      false,
		Sampling:         nil,
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	var err error
	logger, err := config.Build(zap.AddCallerSkip(0))
	if err != nil {
		log.Printf("failed build zap log: %v", err)
		return zap.NewNop()
	}

	zap.RedirectStdLog(logger)

	return logger
}
