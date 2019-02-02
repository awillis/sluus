package processor

import (
	"github.com/awillis/sluus/message"
	"github.com/pkg/errors"
	"runtime"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/plugin"
)

var (
	ErrPluginLoad      = errors.New("unable to load plugin")
	ErrInputClosed     = errors.New("input channel closed")
	ErrBatchProcess    = errors.New("batch process error")
	ErrUncleanShutdown = errors.New("unclean shutdown")
)

type (
	Interface interface {
		ID() string
		Type() plugin.Type
		Plugin() plugin.Interface
		Sluus() *Sluus
		Initialize() (err error)
		Logger() *zap.SugaredLogger
		SetLogger(*zap.SugaredLogger)
		Start()
		Stop()
	}

	Processor struct {
		id         string
		Name       string
		wg         *sync.WaitGroup
		pluginType plugin.Type
		plugin     plugin.Interface
		logger     *zap.SugaredLogger
		sluus      *Sluus
	}
)

func New(name string, pluginType plugin.Type) (proc *Processor) {

	sluus := new(Sluus)

	switch pluginType {
	case plugin.SOURCE:
		sluus.output = make(chan *message.Batch)
	case plugin.CONDUIT:
		sluus.input = make(chan *message.Batch)
		sluus.output = make(chan *message.Batch)
	case plugin.SINK:
		sluus.input = make(chan *message.Batch)
	}

	return &Processor{
		id:         uuid.New().String(),
		Name:       name,
		wg:         new(sync.WaitGroup),
		pluginType: pluginType,
		sluus:      sluus,
	}
}

func (p *Processor) Load() (err error) {
	if plug, e := plugin.New(p.Name, p.pluginType); e != nil {
		err = errors.Wrap(ErrPluginLoad, e.Error())
	} else {
		p.plugin = plug
	}
	return
}

func (p *Processor) Initialize() (err error) {
	p.plugin.SetLogger(p.Logger())
	return p.plugin.Initialize()
}

func (p *Processor) ID() string {
	return p.id
}

func (p *Processor) Type() plugin.Type {
	return p.pluginType
}

func (p *Processor) Plugin() plugin.Interface {
	return p.plugin
}

func (p *Processor) Sluus() *Sluus {
	return p.sluus
}

func (p *Processor) Logger() *zap.SugaredLogger {
	return p.logger.With("name", p.Name, "proc_id", p.ID())
}

func (p *Processor) SetLogger(logger *zap.SugaredLogger) {
	p.logger = logger
}

func (p *Processor) Start() {

	for i := 0; i < runtime.NumCPU(); i++ {
		go func(p *Processor) {
			p.wg.Add(1)

		shutdown:
			for {

				if p.pluginType == plugin.SOURCE {
					if plug, ok := (p.plugin).(plugin.Producer); ok {
						select {
						case output, ok := <-plug.Produce():
							if !ok {
								close(p.Sluus().Output())
								break shutdown
							} else {
								p.Sluus().outputCounter += output.Count()
								p.Sluus().Output() <- output
							}
						}
					}
				}

				if p.pluginType == plugin.CONDUIT {
					select {
					case input, ok := <-p.Sluus().Input():
						if !ok {
							p.Logger().Error(ErrInputClosed)
							close(p.Sluus().Output())
							break shutdown
						} else {
							if plug, ok := (p.plugin).(plugin.Processor); ok {
								p.Sluus().inputCounter += input.Count()
								output, reject, accept, e := plug.Process(input)
								if e != nil {
									p.Logger().Error(errors.Wrap(ErrBatchProcess, e.Error()))
								}
								p.Sluus().outputCounter += output.Count()
								p.Sluus().Output() <- output
								p.Sluus().Reject() <- reject
								p.Sluus().Accept() <- accept
							}
						}
					}
				}

				if p.pluginType == plugin.SINK {
					select {
					case input, ok := <-p.Sluus().Input():
						if !ok {
							p.Logger().Error(ErrInputClosed)
							break shutdown
						} else {
							if plug, ok := (p.plugin).(plugin.Consumer); ok {
								p.Sluus().inputCounter += input.Count()
								plug.Consume() <- input
							}
						}
					}
				}
			}
			p.wg.Done()
		}(p)
	}
}

func (p *Processor) Stop() {
	if plug, ok := (p.plugin).(plugin.Producer); ok {
		if e := plug.Shutdown(); e != nil {
			p.Logger().Error(errors.Wrap(ErrUncleanShutdown, e.Error()))
		}
	}

	if plug, ok := (p.plugin).(plugin.Processor); ok {
		if e := plug.Shutdown(); e != nil {
			p.Logger().Error(errors.Wrap(ErrUncleanShutdown, e.Error()))
		}
	}

	if plug, ok := (p.plugin).(plugin.Consumer); ok {
		if e := plug.Shutdown(); e != nil {
			p.Logger().Error(errors.Wrap(ErrUncleanShutdown, e.Error()))
		}
	}
	p.wg.Wait()
}
