package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"search-service/configs"
)

var Logger *zap.Logger

func NewLogger() *zap.Logger {
	var (
		development   bool
		level         zap.AtomicLevel
		encoderConfig zapcore.EncoderConfig
	)

	level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if configs.Get("SHOW_DEBUG_LOG") == "true" {
		level = zap.NewAtomicLevelAt(zap.DebugLevel)
		development = true
		encoderCfg := zap.NewDevelopmentEncoderConfig()
		encoderConfig = encoderCfg
	}

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderConfig = encoderCfg

	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000 MST")

	config := zap.Config{
		Level:       level,
		Development: development,
		//Sampling is enabled at 100:100 by default,
		// meaning that after the first 100 log entries with the same level and message in the same second,
		// it will log every 100th entry with the same level and message in the same second.
		// You may disable this behavior by setting Sampling to nil.
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    encoderConfig,
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stdout"},
		InitialFields:    nil,
		DisableCaller:    true,
	}
	Logger, _ = config.Build()
	return Logger
}
