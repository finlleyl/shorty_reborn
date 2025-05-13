package logger 

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewSugared() (*zap.SugaredLogger, func(),error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := cfg.Build()
	if err != nil {
		return nil, nil, err
	}

	sugar := logger.Sugar()
	cleanup := func() {
		_ = logger.Sync()
	}

	return sugar, cleanup, nil
}