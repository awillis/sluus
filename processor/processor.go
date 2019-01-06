package processor

import (
	"github.com/golang-collections/go-datastructures/queue"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/plugin"
)

const (
	CONDUIT Category = iota
	SOURCE
	SINK
)

type Category int

type Processor struct {
	id       uuid.UUID
	Name     string
	Logger   *zap.SugaredLogger
	category Category
	plugin   plugin.Processor
	input    chan<- core.Batch
	output   <-chan core.Batch
	queue    *queue.PriorityQueue
}

func NewProcessor(name string, category Category, logger *zap.SugaredLogger) Processor {

	proc := Processor{
		id:       uuid.New(),
		Name:     name,
		category: category,
		input:    make(chan<- core.Batch),
		output:   make(<-chan core.Batch),
		queue:    new(queue.PriorityQueue),
	}

	switch category {
	case CONDUIT:
		proc.plugin = new(Conduit)
	case SOURCE:
		proc.plugin = new(Source)
	case SINK:
		proc.plugin = new(Sink)
	}

	proc.Logger = logger
	proc.plugin.Load(proc.Name)
	return proc
}

func (p Processor) ID() uuid.UUID {
	return p.id
}

func (p Processor) Input() chan<- core.Batch {
	return p.input
}

func (p Processor) Output() <-chan core.Batch {
	return p.output
}
