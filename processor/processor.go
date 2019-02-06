package processor

import (
	"github.com/pkg/errors"
	"runtime"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/plugin"
)

var (
	ErrPluginLoad      = errors.New("unable to load plugin")
	ErrPluginUnknown   = errors.New("unknown plugin type")
	ErrBatchProcess    = errors.New("batch process error")
	ErrUncleanShutdown = errors.New("unclean shutdown")
	ErrProcInterface   = errors.New("processor does not implement interface")
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
	return p.logger.With("processor_id", p.ID())
}

func (p *Processor) SetLogger(logger *zap.SugaredLogger) {
	p.logger = logger
}

func (p *Processor) Start() (err error) {
	for i := 0; i < runtime.NumCPU(); i++ {
		switch p.pluginType {
		case plugin.SOURCE:
			if plug, ok := (p.plugin).(plugin.Producer); ok {
				go startSource(p, plug)
			} else {
				p.Logger().Error(ErrProcInterface)
			}
		case plugin.CONDUIT:
			if plug, ok := (p.plugin).(plugin.Processor); ok {
				go startConduit(p, plug)
			} else {
				p.Logger().Error(ErrProcInterface)
			}
		case plugin.SINK:
			if plug, ok := (p.plugin).(plugin.Consumer); ok {
				go startSink(p, plug)
			} else {
				p.Logger().Error(ErrProcInterface)
			}
		default:
			return ErrPluginUnknown
		}
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
	p.Sluus().Flush()
}

func startSource(p *Processor, plug plugin.Producer) {
	p.wg.Add(1)
	defer p.wg.Done()

shutdown:
	for {
		output, err := plug.Produce()

		if output.Count() > 0 {
			p.Sluus().Output(output)
		}

		if err != nil {
			if err == plugin.ErrShutdown {
				break shutdown
			} else {
				p.Logger().Error(err)
			}
		}
	}
}

func startConduit(p *Processor, plug plugin.Processor) {
	p.wg.Add(1)
	defer p.wg.Done()

shutdown:
	for {
		input := p.Sluus().Input()
		output, reject, accept, err := plug.Process(input)

		if output.Count() > 0 {
			p.Sluus().Output(output)
		}

		if reject.Count() > 0 {
			p.Sluus().Reject(reject)
		}

		if accept.Count() > 0 {
			p.Sluus().Accept(accept)
		}

		if err != nil {
			if err == plugin.ErrShutdown {
				break shutdown
			} else {
				p.Logger().Error(err)
			}
		}
	}
}

func startSink(p *Processor, plug plugin.Consumer) {
	p.wg.Add(1)
	defer p.wg.Done()

shutdown:
	for {
		input, err := p.Sluus().Input()
		plug.Consume(input)

		if err != nil {
			if err == plugin.ErrShutdown {
				break shutdown
			} else {
				p.Logger().Error(err)
			}
		}
	}
}
