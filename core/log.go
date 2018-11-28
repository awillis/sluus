package core

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger zap.SugaredLogger

func NewLogger(component string) *zap.SugaredLogger {

	cfg := zap.Config{
		Level:             zap.AtomicLevel{},
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: false,
		Sampling: &zap.SamplingConfig{
			Initial:    0,
			Thereafter: 0,
		},
		Encoding: "",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "",
			LevelKey:       "",
			TimeKey:        "",
			NameKey:        "",
			CallerKey:      "",
			StacktraceKey:  "",
			LineEnding:     "",
			EncodeLevel:    nil,
			EncodeTime:     nil,
			EncodeDuration: nil,
			EncodeCaller:   nil,
			EncodeName:     nil,
		},
		OutputPaths:      nil,
		ErrorOutputPaths: nil,
		InitialFields:    nil,
	}

	_ = cfg
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}
