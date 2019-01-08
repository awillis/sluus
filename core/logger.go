package core

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func SetupLogger(logfile string) *zap.SugaredLogger {
	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("unable to instantiate logger: %v", err)
	}

	output := zapcore.Lock(f)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	zcore := zapcore.NewTee(zapcore.NewCore(encoder, output, priority))
	return zap.New(zcore).Sugar()
}
