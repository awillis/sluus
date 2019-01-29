package processor

import (
	"context"
	"github.com/pkg/errors"
	"runtime"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/plugin"
)

var ErrPluginLoad = errors.New("unable to load plugin")

type (
	Interface interface {
		ID() string
		Type() plugin.Type
		Options() interface{}
		Sluus() *Sluus
		Logger() *zap.SugaredLogger
		SetLogger(*zap.SugaredLogger)
		Run()
	}

	Processor struct {
		id         string
		Name       string
		wg         *sync.WaitGroup
		pluginType plugin.Type
		plugin     plugin.Processor
		context    context.Context
		logger     *zap.SugaredLogger
		sluus      *Sluus
	}

	ContextKey struct {
	}
)

func New(name string, pluginType plugin.Type) (proc *Processor) {

	return &Processor{
		id:         uuid.New().String(),
		Name:       name,
		wg:         new(sync.WaitGroup),
		pluginType: pluginType,
		context:    context.Background(),
		flume:      new(Flume),
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

func (p *Processor) Run() {

	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			p.wg.Add(1)
		shutdown:
			for {
				select {
				case input, ok := <-p.Flume().Output():
					if !ok {
						p.Logger().Error("output channel closed")
						break shutdown
					} else {
						if p.pluginType != plugin.SOURCE {
							pass, reject, accept, e := p.plugin.Process(input)
							if e != nil {
								p.Logger().Error(e)
							}
							p.Flume().Input() <- pass
							p.Flume().Input() <- reject
							p.Flume().Input() <- accept
						}
					}
				}
			}
			p.wg.Done()
		}()
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

func (p *Processor) Flume() *Flume {
	return p.flume
}

func (p *Processor) Logger() *zap.SugaredLogger {
	return p.logger.With("processor", p.ID())
}

func (p *Processor) SetLogger(logger *zap.SugaredLogger) {
	p.logger = logger
}
