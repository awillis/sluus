package core

import (
	"os"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func SetupLogger(conf *zap.Config) *zap.SugaredLogger {

	logger, err := conf.Build()
	if err != nil {
		panic(err)
	}
	return logger.Sugar()
}

func LogConfig(name string, id string) *zap.Config {

	basename := name
	fields := make(map[string]interface{})

	if name == "core" {
		fields["pid"] = id
	} else {
		basename = name + "-" + id
		fields["pipe"] = name
		fields["pipe_id"] = id
	}

	return &zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		DisableStacktrace: true,
		DisableCaller:     false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:  "msg",
			LevelKey:    "level",
			TimeKey:     "time",
			LineEnding:  zapcore.DefaultLineEnding,
			EncodeLevel: zapcore.LowercaseLevelEncoder,
			EncodeTime: zapcore.TimeEncoder(func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
				enc.AppendString(t.Format(time.RFC3339))
			}),
			EncodeDuration: zapcore.SecondsDurationEncoder,
		},
		OutputPaths:   []string{strings.Join([]string{LOGDIR, basename + ".log"}, string(os.PathSeparator))},
		InitialFields: fields,
	}
}
