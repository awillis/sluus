package core

import (
	"net/url"
	"os"
	"runtime"
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
	}

	logfile := new(strings.Builder)

	if runtime.GOOS == "windows" {
		logfile.WriteString("windows:///")

		// fix for windows paths: https://github.com/uber-go/zap/issues/621
		_ = zap.RegisterSink("windows", func(i *url.URL) (sink zap.Sink, e error) {
			return os.OpenFile(i.Path[1:], os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		})
	}

	logfile.WriteString(LOGDIR)
	logfile.WriteRune(os.PathSeparator)
	logfile.WriteString(basename)
	logfile.WriteString(".log")

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
		OutputPaths:   []string{logfile.String()},
		InitialFields: fields,
	}
}
