package pipeline

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"os"
	"strings"
	"time"

	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor"
)

var (
	ErrInvalidProcessor = errors.New("invalid processor")
	ErrNoSource         = errors.New("missing source processor")
)

type (
	Component struct {
		id   uint
		next *Component
		pipe *Pipe
		proc *processor.Processor
	}

	Pipe struct {
		Id        string
		Name      string
		logger    *zap.SugaredLogger
		hasSource bool
		hasAccept bool
		hasReject bool
		root      Component
		reject    *Component
		accept    *Component
		len       uint
	}
)

func (c *Component) Next() (next *Component) {
	if c.pipe != nil && c.next != &c.pipe.root {
		next = c.next
	}
	return
}

func New(name string) (pipe *Pipe) {

	pipe = new(Pipe)
	pipe.Name = name
	pipe.Id = uuid.New().String()
	pipe.root.next = &pipe.root
	pipe.len = 0

	// Setup logger
	pipe.logger = core.SetupLogger(core.LogConfig(pipe.Name, pipe.ID()))
	return pipe
}

func (p *Pipe) ID() string {
	return p.Id
}

func (p *Pipe) Logger() *zap.SugaredLogger {
	return p.logger
}

func (p *Pipe) Len() uint {
	return p.len
}

func (p *Pipe) Source() *Component {
	if p.len == 0 {
		return nil
	}
	return p.root.next
}

func (p *Pipe) Accept() *Component {
	if p.len == 0 {
		return nil
	}
	return p.accept
}

func (p *Pipe) Reject() *Component {
	if p.len == 0 {
		return nil
	}
	return p.reject
}

func (p *Pipe) Start(ctx context.Context) {
	for n := p.Source(); n != nil; n = n.Next() {
		if err := n.proc.Start(ctx); err != nil {
			p.Logger().Error(err)
		}
	}
}

func (p *Pipe) Stop() {
	for n := p.Source(); n != nil; n = n.Next() {
		n.proc.Stop()
	}
}

func (p *Pipe) Attach(component *Component) {

	for n := &p.root; n != component; n = n.Next() {
		if n.Next() == nil {
			component.proc.SetLogger(p.logger)
			p.len++
			component.id = p.len
			n.pipe = p
			n.next = component
		}
	}
}

func (p *Pipe) ConfigureAndInitialize(pipeConf PipeConfig) {

	pollIntvl := processor.PollInterval(time.Duration(pipeConf.PollInterval) * time.Millisecond)
	batchSize := processor.BatchSize(pipeConf.BatchSize)
	batchTimeout := processor.BatchTimeout(time.Duration(pipeConf.BatchTimeout) * time.Second)
	ringSize := processor.RingSize(pipeConf.RingSize)
	queueinFlight := processor.QueryInFlight(pipeConf.QueueDepth)
	tableMode := processor.TableLoadingMode(pipeConf.TableLoadingMode)
	valueMode := processor.ValueLogLoadingMode(pipeConf.ValueLogLoadingMode)

	for n := p.Source(); n != nil; n = n.Next() {

		dir := dataDirBuilder(p.Name)

		dir.WriteString(fmt.Sprintf("%d-%s-%s",
			n.id,
			plugin.TypeName(n.proc.Plugin().Type()),
			n.proc.Plugin().Name()))

		dataDir := processor.DataDir(dir.String())

		if err := n.proc.Configure(
			ringSize,
			queueinFlight,
			pollIntvl,
			batchSize,
			batchTimeout,
			dataDir,
			tableMode,
			valueMode,
		); err != nil {
			p.Logger().Error(err)
		}

		if e := n.proc.Sluus().Initialize(); e != nil {
			p.Logger().Error(e)
		}

		if n != p.Accept() {
			input := processor.Input(n.proc.Sluus().Output())
			if err := n.Next().proc.Configure(input); err != nil {
				p.Logger().Error(err)
			}
		}
	}

	for n := p.Source(); n != p.Reject(); n = n.Next() {

		reject := processor.Reject(p.Reject().proc.Sluus().Input())
		accept := processor.Accept(p.Accept().proc.Sluus().Input())

		if err := n.proc.Configure(
			reject,
			accept,
		); err != nil {
			p.Logger().Error(err)
		}
		if e := n.proc.Initialize(); e != nil {
			p.Logger().Error(e)
		}
	}
}

func (p *Pipe) Add(proc *processor.Processor) (err error) {

	component := new(Component)
	component.proc = proc

	switch proc.Type() {
	case plugin.SOURCE:
		if p.hasSource {
			return ErrInvalidProcessor
		} else {
			p.hasSource = true
		}
	case plugin.CONDUIT:
		if !p.hasSource {
			return ErrNoSource
		}
	case plugin.SINK:
		if !p.hasSource {
			return ErrNoSource
		}
		if !p.hasReject {
			p.reject = component
			p.hasReject = true
		} else if !p.hasAccept {
			p.accept = component
			p.hasAccept = true
		} else {
			return ErrInvalidProcessor
		}
	default:
		return ErrInvalidProcessor
	}

	p.Attach(component)
	return
}

func dataDirBuilder(pipeName string) (dirpath *strings.Builder) {
	dirpath = new(strings.Builder)
	dirpath.WriteString(core.DATADIR)
	dirpath.WriteRune(os.PathSeparator)
	dirpath.WriteString(pipeName)
	dirpath.WriteRune(os.PathSeparator)
	return
}
