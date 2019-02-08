package processor

import (
	"github.com/awillis/sluus/message"
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

	runner struct {
		// plugin interface
		produce func() (*message.Batch, error)
		process func(*message.Batch) (output, reject, accept *message.Batch, err error)
		consume func(*message.Batch) error
		logger  func(...interface{})

		// sluus
		receive  func() *message.Batch
		output   func(*message.Batch)
		reject   func(*message.Batch)
		accept   func(*message.Batch)
		shutdown func()
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

		// runner is used to avoid interface dynamic dispatch penalty
		runner := new(runner)
		runner.logger = p.logger.Error
		runner.receive = p.sluus.receive
		runner.output = p.sluus.sendOutput
		runner.reject = p.sluus.sendReject
		runner.accept = p.sluus.sendAccept
		runner.shutdown = p.sluus.shutdown

		switch p.pluginType {
		case plugin.SOURCE:
			if plug, ok := (p.plugin).(plugin.Producer); ok {
				runner.produce = plug.Produce
				go runSource(p, runner)
			} else {
				p.Logger().Error(ErrProcInterface)
			}
		case plugin.CONDUIT:
			if plug, ok := (p.plugin).(plugin.Processor); ok {
				runner.process = plug.Process
				go runConduit(p, runner)
			} else {
				p.Logger().Error(ErrProcInterface)
			}
		case plugin.SINK:
			if plug, ok := (p.plugin).(plugin.Consumer); ok {
				runner.consume = plug.Consume
				go runSink(p, runner)
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

func runSource(p *Processor, r *runner) {
	p.wg.Add(1)
	defer p.wg.Done()
	defer r.shutdown()

	for {
		output, err := r.produce()

		if err != nil {
			if err == plugin.ErrShutdown {
				break
			}
			r.logger(err)
		} else {
			r.output(output)
		}
	}
}

func runConduit(p *Processor, r *runner) {
	p.wg.Add(1)
	defer p.wg.Done()
	defer r.shutdown()

	for {
		input := r.receive()
		output, reject, accept, err := r.process(input)
		r.output(output)
		r.reject(reject)
		r.accept(accept)

		if err == plugin.ErrShutdown {
			break
		}
	}
}

func runSink(p *Processor, r *runner) {
	p.wg.Add(1)
	defer p.wg.Done()
	defer r.shutdown()

	for {
		input := r.receive()
		if err := r.consume(input); err != nil {
			if err == plugin.ErrShutdown {
				break
			}
			r.logger(err)
		}
	}
}
