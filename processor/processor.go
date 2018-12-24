package processor

import (
	"github.com/golang-collections/go-datastructures/queue"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
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
	plugin   plugin.Component
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
		proc.plugin = new(Conduit)
	case SOURCE:
		proc.plugin = new(Source)
		close(proc.input)
	case SINK:
		proc.plugin = new(Sink)
		close(proc.output)
	}

	proc.Logger = logger
	proc.plugin.Load(proc.Name)
	proc.plugin.(*Source).Produce()
	proc.plugin.(*Conduit).Convey()
	proc.plugin.(*Sink).Consume()
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
