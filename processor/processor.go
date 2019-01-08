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
	id       uuid.UUID
	Name     string
	Logger   *zap.SugaredLogger
	Context  context.Context
	plugtype core.PluginType
	plugin   plugin.Processor
	input    chan<- core.Batch
	output   <-chan core.Batch
	queue    *queue.PriorityQueue
}

func NewProcessor(pluginName string, pluginType core.PluginType, logger *zap.SugaredLogger) Processor {

	proc := Processor{
		id:       uuid.New(),
		Name:     pluginName,
		plugtype: pluginType,
		input:    make(chan<- core.Batch),
		output:   make(<-chan core.Batch),
		queue:    new(queue.PriorityQueue),
	}

	proc.Logger = logger
	plug, err := plugin.Load(pluginName, pluginType)
	if err != nil {
		proc.Logger.Errorf("unable to load plugin: %s: %s", pluginName, err)
	}

	proc.plugin = plug
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
