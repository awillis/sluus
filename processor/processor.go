package processor

import (
	"context"

	"github.com/golang-collections/go-datastructures/queue"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

type Interface interface {
	ID() string
	Type() plugin.Type
	Input() chan<- message.Batch
	Output() <-chan message.Batch
	SetLogger(logger *zap.SugaredLogger)
	Logger() *zap.SugaredLogger
}

type Processor struct {
	id      string
	Name    string
	logger  *zap.SugaredLogger
	Context context.Context
	ptype   plugin.Type
	plugin  plugin.Processor
	input   chan<- message.Batch
	output  <-chan message.Batch
	queue   *queue.PriorityQueue
}

func NewProcessor(name string, ptype plugin.Type) *Processor {

	proc := &Processor{
		id:     uuid.New().String(),
		Name:   name,
		ptype:  ptype,
		input:  make(chan<- message.Batch),
		output: make(<-chan message.Batch),
		queue:  new(queue.PriorityQueue),
	}

	plug, err := plugin.Load(name, ptype)
	if err != nil {
		core.Logger.Errorf("unable to load plugin: %s: %s", name, err)
	}

	proc.plugin = plug
	return proc
}

func (p Processor) ID() string {
	return p.id
}

func (p Processor) Type() plugin.Type {
	return p.plugin.Type()
}

func (p Processor) Input() chan<- message.Batch {
	return p.input
}

func (p Processor) Output() <-chan message.Batch {
	return p.output
}

func (p Processor) SetLogger(logger *zap.SugaredLogger) {
	p.logger = logger
}

func (p Processor) Logger() *zap.SugaredLogger {
	return p.logger.With("processor", p.ID())
}
