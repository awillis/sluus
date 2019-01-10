package processor

import (
	"context"

	"github.com/golang-collections/go-datastructures/queue"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/plugin"
)

type Interface interface {
	ID() string
	Type() plugin.Type
	Input() chan<- core.Batch
	Output() <-chan core.Batch
	SetLogger(logger *zap.SugaredLogger)
}

type Processor struct {
	id      string
	Name    string
	Logger  *zap.SugaredLogger
	Context context.Context
	ptype   plugin.Type
	plugin  plugin.Interface
	input   chan<- core.Batch
	output  <-chan core.Batch
	queue   *queue.PriorityQueue
}

func NewProcessor(name string, ptype plugin.Type) *Processor {

	proc := &Processor{
		id:     uuid.New().String(),
		Name:   name,
		ptype:  ptype,
		input:  make(chan<- core.Batch),
		output: make(<-chan core.Batch),
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

func (p Processor) Input() chan<- core.Batch {
	return p.input
}

func (p Processor) Output() <-chan core.Batch {
	return p.output
}

func (p Processor) SetLogger(logger *zap.SugaredLogger) {
	p.Logger = logger
}
