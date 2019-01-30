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

var (
	ErrPluginLoad   = errors.New("unable to load plugin")
	ErrInputClosed  = errors.New("input channel closed")
	ErrBatchProcess = errors.New("batch process error")
)

type (
	Interface interface {
		ID() string
		Type() plugin.Type
		Plugin() plugin.Loader
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
		plugin     plugin.Loader
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
		sluus:      new(Sluus),
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
		go func(p *Processor) {
			p.wg.Add(1)
		shutdown:
			for {
				if p.pluginType == plugin.SOURCE {
					if plug, ok := (p.plugin).(plugin.Producer); ok {
						if err := plug.Produce(); err != nil {

						}
					} else {
						select {
						case input, ok := <-p.Sluus().Output():
							if !ok {
								p.Logger().Error(ErrInputClosed)
								break shutdown
							} else {
								if plug, ok := (p.plugin).(plugin.Processor); ok {
									p.Sluus().inputCounter += input.Count()
									pass, reject, accept, e := plug.Process(input)
									if e != nil {
										p.Logger().Error(errors.Wrap(ErrBatchProcess, e.Error()))
									}
									p.Sluus().outputCounter += pass.Count()
									p.Sluus().Output() <- pass
									p.Sluus().Reject() <- reject
									p.Sluus().Accept() <- accept
								}
							}
						}
					}
				}
			}
			p.wg.Done()
		}(p)
	}
}

func (p *Processor) ID() string {
	return p.id
}

func (p *Processor) Type() plugin.Type {
	return p.pluginType
}

func (p *Processor) Plugin() plugin.Loader {
	return p.plugin
}

func (p *Processor) Sluus() *Sluus {
	return p.sluus
}

func (p *Processor) Logger() *zap.SugaredLogger {
	return p.logger.With("processor", p.ID())
}

func (p *Processor) SetLogger(logger *zap.SugaredLogger) {
	p.logger = logger
}
