package processor

import (
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
		logger     *zap.SugaredLogger
		input      chan<- message.Batch
		output     <-chan message.Batch
	}
)

func New(name string, pluginType plugin.Type) (proc *Processor) {

	return &Processor{
		id:         uuid.New().String(),
		Name:       name,
		pluginType: pluginType,
		input:      make(chan<- message.Batch),
		output:     make(<-chan message.Batch),
	}
}

func (p *Processor) Load() (err error) {
	if plug, err := plugin.NewProcessor(p.Name, p.pluginType); err != nil {
		return errors.Wrapf(err, ErrPluginLoad.Error())
	} else {
		p.plugin = plug
	}

	return
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
