package logger

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Mode string

const (
	Local Mode = "local"
	Dev   Mode = "dev"
	Prod  Mode = "prod"
)

func NewSugared(mode Mode) (*zap.SugaredLogger, func(), error) {
	var cfg zap.Config

	switch mode {
	case Local:
		cfg = zap.NewDevelopmentConfig()
		cfg.Encoding = "console"
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	case Dev:
		cfg = zap.NewDevelopmentConfig()
		cfg.Encoding = "json"
		cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	case Prod:
		cfg = zap.NewProductionConfig()
		cfg.Encoding = "json"

	default:
		return nil, nil, fmt.Errorf("invalid mode: %q", mode)
	}

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
