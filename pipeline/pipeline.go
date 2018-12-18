package pipeline

import (
	"fmt"
	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/processor"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

type Pipeline interface {
	ID() uuid.UUID
	Logger() *zap.SugaredLogger
}

type Pipe struct {
	id     uuid.UUID
	logger *zap.SugaredLogger
}

func NewPipeline() *Pipe {

	pipe := new(Pipe)
	pipe.id = uuid.New()

	// Setup logger
	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	})

	logfile := strings.Join([]string{core.LOGDIR, "pipeline_" + pipe.ID().String()}, string(os.PathSeparator))
	f, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("unable to instantiate pipeline logger: %v", err)
	}

	output := zapcore.Lock(f)
	encoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	zcore := zapcore.NewTee(zapcore.NewCore(encoder, output, priority))
	pipe.logger = zap.New(zcore).Sugar()
	return pipe
}

func (p *Pipe) ID() uuid.UUID {
	return p.id
}

func (p *Pipe) Logger() *zap.SugaredLogger {
	return p.logger
}

func (p *Pipe) AddProcessor(name string, category processor.Category) {
	proc := processor.NewProcessor(name, category, p.logger)
	proc.ID()
}
