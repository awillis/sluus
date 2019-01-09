package processor

import (
	"context"

	"github.com/golang-collections/go-datastructures/queue"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/plugin"
)

type Processor struct {
	id       string
	Name     string
	Logger   *zap.SugaredLogger
	Context  context.Context
	plugtype plugin.Type
	plugin   plugin.Interface
	input    chan<- core.Batch
	output   <-chan core.Batch
	queue    *queue.PriorityQueue
}

func NewProcessor(name string, ptype plugin.Type, logger *zap.SugaredLogger) Processor {

	proc := Processor{
		id:       uuid.New().String(),
		Name:     name,
		plugtype: ptype,
		input:    make(chan<- core.Batch),
		output:   make(<-chan core.Batch),
		queue:    new(queue.PriorityQueue),
	}

	proc.Logger = logger
	plug, err := plugin.Load(name, ptype)
	if err != nil {
		proc.Logger.Errorf("unable to load plugin: %s: %s", name, err)
	}

	proc.plugin = plug
	return proc
}

func (p Processor) ID() string {
	return p.id
}

func (p Processor) Input() chan<- core.Batch {
	return p.input
}

func (p Processor) Output() <-chan core.Batch {
	return p.output
}
