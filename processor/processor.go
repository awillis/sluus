package processor

import (
	"context"
	//"runtime"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/awillis/sluus/message"
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
	Processor struct {
		id           string
		Name         string
		wg           *sync.WaitGroup
		pluginType   plugin.Type
		plugin       plugin.Interface
		pollInterval time.Duration
		logger       *zap.SugaredLogger
		sluus        *Sluus
	}

	runner struct {
		// plugin interface
		start   func(context.Context)
		produce func() <-chan *message.Batch
		process func(*message.Batch) *message.Batch
		consume func(*message.Batch)

		// sluus
		receive func() <-chan *message.Batch
		output  func(*message.Batch)
		reject  func(*message.Batch)
		accept  func(*message.Batch)
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
	p.sluus = newSluus(p.pluginType)

	if plug, e := plugin.New(p.Name, p.pluginType); e != nil {
		err = errors.Wrap(ErrPluginLoad, e.Error())
	} else {
		p.plugin = plug
	}
	return
}

func (p *Processor) Initialize() (err error) {
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
	return p.logger.With("processor_id", p.ID(), "name", p.Name, "type", plugin.TypeName(p.pluginType))
}

func (p *Processor) SetLogger(logger *zap.SugaredLogger) {
	p.logger = logger
	p.sluus.SetLogger(logger)
	p.plugin.SetLogger(logger)
}

func (p *Processor) Start(ctx context.Context) (err error) {

	p.sluus.Start()

	// runner is used to avoid interface dynamic dispatch penalty
	runner := new(runner)
	runner.receive = p.sluus.receiveInput
	runner.output = p.sluus.sendOutput
	runner.reject = p.sluus.sendReject
	runner.accept = p.sluus.sendAccept

	switch p.pluginType {
	case plugin.SOURCE:
		if plug, ok := (p.plugin).(plugin.Producer); ok {
			runner.start = plug.Start
			runner.produce = plug.Produce
			go runSource(p, ctx, runner)
		} else {
			p.Logger().Error(ErrProcInterface)
		}
	case plugin.CONDUIT:
		if plug, ok := (p.plugin).(plugin.Processor); ok {
			runner.start = plug.Start
			runner.process = plug.Process
			go runConduit(p, ctx, runner)
		} else {
			p.Logger().Error(ErrProcInterface)
		}
	case plugin.SINK:
		if plug, ok := (p.plugin).(plugin.Consumer); ok {
			runner.start = plug.Start
			runner.consume = plug.Consume
			go runSink(p, ctx, runner)
		} else {
			p.Logger().Error(ErrProcInterface)
		}
	default:
		return ErrPluginUnknown
	}
	return
}

func runSource(p *Processor, ctx context.Context, r *runner) {
	p.wg.Add(1)
	defer p.wg.Done()
	r.start(ctx)

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

loop:
	select {
	case <-ctx.Done():
		break
	case <-ticker.C:
		// runtime.Gosched()
		goto loop
	case batch, ok := <-r.produce():
		if ok {
			p.sluus.outCtr += batch.Count()
			r.output(batch)
		}
		// runtime.Gosched()
		goto loop
	}
}

func runConduit(p *Processor, ctx context.Context, r *runner) {
	p.wg.Add(1)
	defer p.wg.Done()
	r.start(ctx)

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

loop:
	select {
	case <-ctx.Done():
		break
	case <-ticker.C:
		// runtime.Gosched()
		goto loop
	case batch, ok := <-r.receive():
		if ok {
			pBatch := r.process(batch)
			r.output(pBatch.Pass())
			r.reject(pBatch.Reject())
			r.accept(pBatch.Accept())
		}
		// runtime.Gosched()
		goto loop
	}
}

func runSink(p *Processor, ctx context.Context, r *runner) {

	p.wg.Add(1)
	defer p.wg.Done()

	r.start(ctx)

	ticker := time.NewTicker(p.pollInterval)
	defer ticker.Stop()

loop:
	select {
	case <-ctx.Done():
		break
	case <-ticker.C:
		// runtime.Gosched()
		goto loop
	case batch, ok := <-r.receive():
		if ok {
			p.sluus.inCtr += batch.Count()
			r.consume(batch)
		}
		// runtime.Gosched()
		goto loop
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
	p.Logger().Info("processor wait")
	p.wg.Wait()
	p.sluus.shutdown()
	p.Logger().Info("processor shutdown")
}
