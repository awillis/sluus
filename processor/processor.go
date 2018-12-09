package processor

import (
	"github.com/golang-collections/go-datastructures/queue"
	"github.com/google/uuid"

	"sluus/core"
	"sluus/plugin"
	"sluus/processor/conduit"
	"sluus/processor/sink"
	"sluus/processor/source"
)

const (
	SOURCE Category = iota + 1
	CONDUIT
	SINK
)

type Category int

type Processor interface {
	ID() uuid.UUID
	Input() chan core.Message
	Output() chan core.Message
	Process()
}

type Base struct {
	id       uuid.UUID
	Name     string
	category Category
	plugin   plugin.Plugin
	input    chan core.Message
	output   chan core.Message
	queue    *queue.PriorityQueue
}

func NewProcessor(name string, category Category) Base {

	proc := Base{
		id:       uuid.New(),
		Name:     name,
		category: category,
		input:    make(chan core.Message),
		output:   make(chan core.Message),
		queue:    new(queue.PriorityQueue),
	}

	switch category {
	case SOURCE:
		proc.plugin = new(source.Source)
	case CONDUIT:
		proc.plugin = new(conduit.Conduit)
	case SINK:
		proc.plugin = new(sink.Sink)
	}

	proc.plugin.PluginLoad(proc.Name)
	proc.plugin.(*source.Source).Test()
	return proc
}

func (p *Base) ID() uuid.UUID {
	return p.id
}

func (p *Base) Input() chan core.Message {
	return p.input
}

func (p *Base) Output() chan core.Message {
	return p.output
}
