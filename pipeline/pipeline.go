package pipeline

import (
	"os"
	"strings"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/processor"
)

type Pipeline interface {
	ID() uuid.UUID
	Logger() *zap.SugaredLogger
}

type Pipe struct {
	id         uuid.UUID
	logger     *zap.SugaredLogger
	processors map[string]processor.Processor
	gates      map[string]Sluus
}

func NewPipeline() *Pipe {

	pipe := new(Pipe)
	pipe.id = uuid.New()

	// Setup logger
	logfile := strings.Join([]string{core.LOGDIR, "pipeline_" + pipe.ID().String()}, string(os.PathSeparator))
	pipe.logger = core.SetupLogger(logfile)
	return pipe
}

func (p *Pipe) ID() uuid.UUID {
	return p.id
}

func (p *Pipe) Logger() *zap.SugaredLogger {
	return p.logger
}

func (p *Pipe) AddProcessor(name string, ptype core.PluginType) {
	proc := processor.NewProcessor(name, ptype, p.logger)
	p.processors[proc.ID().String()] = proc
}
