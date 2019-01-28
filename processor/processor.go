package processor

import (
	"context"
	"github.com/pkg/errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

var ErrPluginLoad = errors.New("unable to load plugin")

type (
	Interface interface {
		ID() string
		Type() plugin.Type
		Options() interface{}
		Input() chan<- message.Batch
		Output() <-chan message.Batch
		Logger() *zap.SugaredLogger
		SetLogger(*zap.SugaredLogger)
	}

	Processor struct {
		id         string
		Name       string
		pluginType plugin.Type
		plugin     plugin.Processor
		context    context.Context
		logger     *zap.SugaredLogger
		input      chan<- message.Batch
		output     <-chan message.Batch
	}

	ContextKey struct {
	}
)

func New(name string, pluginType plugin.Type) (proc *Processor) {

	return &Processor{
		id:         uuid.New().String(),
		Name:       name,
		pluginType: pluginType,
		context:    context.Background(),
		input:      make(chan<- message.Batch),
		output:     make(<-chan message.Batch),
	}
}

func (p *Processor) Context() (ctx context.Context) {
	key := new(ContextKey)
	return context.WithValue(p.context, key, p.id)
}

func (p *Processor) Load() (err error) {
	if plug, e := plugin.NewProcessor(p.Name, p.pluginType); e != nil {
		return errors.Wrap(ErrPluginLoad, e.Error())
	} else {
		p.plugin = plug
		return p.plugin.Initialize(p.Context())
	}
}

func (p *Processor) ID() string {
	return p.id
}

func (p *Processor) Type() plugin.Type {
	return p.pluginType
}

func (p *Processor) Options() interface{} {
	return p.plugin.Options()
}

func (p *Processor) Input() chan<- message.Batch {
	return p.input
}

func (p *Processor) Output() <-chan message.Batch {
	return p.output
}

func (p *Processor) Logger() *zap.SugaredLogger {
	return p.logger.With("processor", p.ID())
}

func (p *Processor) SetLogger(logger *zap.SugaredLogger) {
	p.logger = logger
}
