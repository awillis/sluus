package processor

import (
	"context"

	"github.com/golang-collections/go-datastructures/queue"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

type (
	Interface interface {
		ID() string
		Type() plugin.Type
		Input() chan<- message.Batch
		Output() <-chan message.Batch
		SetLogger(logger *zap.SugaredLogger)
		Logger() *zap.SugaredLogger
	}

	Processor struct {
		id         string
		Name       string
		logger     *zap.SugaredLogger
		Context    context.Context
		pluginType plugin.Type
		plugin     plugin.Processor
		input      chan<- message.Batch
		output     <-chan message.Batch
		queue      *queue.PriorityQueue
	}
)

func New(name string, pluginType plugin.Type) (proc *Processor) {

	return &Processor{
		id:         uuid.New().String(),
		Name:       name,
		pluginType: pluginType,
		input:      make(chan<- message.Batch),
		output:     make(<-chan message.Batch),
		queue:      new(queue.PriorityQueue),
	}
}

func (p *Processor) Load(options map[string]interface{}) (err error) {
	if plug, err := plugin.NewProcessor(p.Name, p.pluginType); err != nil {
		p.Logger().Errorf("unable to load plugin: %s: %s", p.Name, err)
	} else {
		p.plugin = plug
	}

	//if err = p.plugin.Initialize(options); err != nil {
	//	return
	//}
	return
}

func (p Processor) ID() string {
	return p.id
}

func (p Processor) Type() plugin.Type {
	return p.pluginType
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
