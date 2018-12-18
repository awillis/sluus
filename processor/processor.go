package processor

import (
	"github.com/golang-collections/go-datastructures/queue"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor/conduit"
	"github.com/awillis/sluus/processor/sink"
	"github.com/awillis/sluus/processor/source"
)

const (
	CONDUIT Category = iota
	SOURCE
	SINK
)

type Category int

type Processor interface {
	ID() uuid.UUID
	Input() chan message.Message
	Output() chan message.Message
}

type Base struct {
	id       uuid.UUID
	Name     string
	Logger   *zap.SugaredLogger
	category Category
	plugin   plugin.Plugin
	input    chan message.Message
	output   chan message.Message
	queue    *queue.PriorityQueue
}

func NewProcessor(name string, category Category, logger *zap.SugaredLogger) Base {

	proc := Base{
		id:       uuid.New(),
		Name:     name,
		category: category,
		input:    make(chan message.Message),
		output:   make(chan message.Message),
		queue:    new(queue.PriorityQueue),
	}

	switch category {
	case CONDUIT:
		proc.plugin = new(conduit.Conduit)
	case SOURCE:
		proc.plugin = new(source.Source)
		close(proc.input)
	case SINK:
		proc.plugin = new(sink.Sink)
		close(proc.output)
	}

	proc.Logger = logger
	proc.plugin.Load(proc.Name)
	proc.plugin.(*source.Source).Produce()
	proc.plugin.(*conduit.Conduit).Convey()
	proc.plugin.(*sink.Sink).Consume()
	return proc
}

func (p Base) ID() uuid.UUID {
	return p.id
}

func (p Base) Input() chan message.Message {
	return p.input
}

func (p Base) Output() chan message.Message {
	return p.output
}
