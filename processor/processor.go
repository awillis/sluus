package processor

import (
	"runtime"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/awillis/sluus/plugin"
)

var (
	ErrPluginLoad      = errors.New("unable to load plugin")
	ErrPluginUnknown   = errors.New("unknown plugin type")
	ErrUncleanShutdown = errors.New("unclean shutdown")
	ErrProcInterface   = errors.New("processor does not implement interface")
	ErrInitialize      = errors.New("initialization err")
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
		Start() error
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
	return &Processor{
		id:         uuid.New().String(),
		Name:       name,
		wg:         new(sync.WaitGroup),
		pluginType: pluginType,
	}
}

func (p *Processor) Load() (err error) {
	p.sluus = NewSluus(p)

	if plug, e := plugin.New(p.Name, p.pluginType); e != nil {
		err = errors.Wrap(ErrPluginLoad, e.Error())
	} else {
		p.plugin = plug
	}
	return
}

func (p *Processor) Initialize() (err error) {
	p.sluus.SetLogger(p.Logger())
	p.plugin.SetLogger(p.Logger())

	if e := p.sluus.Initialize(); e != nil {
		return errors.Wrap(ErrInitialize, e.Error())
	}

	if e := p.plugin.Initialize(); e != nil {
		return errors.Wrap(ErrInitialize, e.Error())
	}
	return
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
	return
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

func startSource(p *Processor, plug plugin.Producer) {
	p.wg.Add(1)
	defer p.wg.Done()
	defer p.Sluus().Shutdown()

	for {
		output, err := plug.Produce()

		if err != nil {
			if err == plugin.ErrShutdown {
				break
			}
			p.Logger().Error(err)
		} else {
			p.Sluus().Pass(output)
		}
	}
}

func startConduit(p *Processor, plug plugin.Processor) {
	p.wg.Add(1)
	defer p.wg.Done()
	defer p.Sluus().Shutdown()

	for {
		input := p.Sluus().Receive()
		output, reject, accept, err := plug.Process(input)
		p.Sluus().Pass(output)
		p.Sluus().Reject(reject)
		p.Sluus().Accept(accept)

		if err == plugin.ErrShutdown {
			break
		}
	}
}

func startSink(p *Processor, plug plugin.Consumer) {
	p.wg.Add(1)
	defer p.wg.Done()
	defer p.Sluus().Shutdown()

	for {
		input := p.Sluus().Receive()
		if err := plug.Consume(input); err != nil {
			if err == plugin.ErrShutdown {
				break
			}
			p.Logger().Error(err)
		}
	}
}
