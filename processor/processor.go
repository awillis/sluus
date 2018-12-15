package processor

import (
	"github.com/golang-collections/go-datastructures/queue"
	"github.com/google/uuid"

	"github.com/awillis/sluus/core"
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
	Input() chan core.Message
	Output() chan core.Message
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
	case CONDUIT:
		proc.plugin = new(conduit.Conduit)
	case SOURCE:
		proc.plugin = new(source.Source)
		close(proc.input)
	case SINK:
		proc.plugin = new(sink.Sink)
		close(proc.output)
	}

	proc.plugin.Load(proc.Name)
	proc.plugin.(*source.Source).Produce()
	proc.plugin.(*conduit.Conduit).Convey()
	proc.plugin.(*sink.Sink).Consume()
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
