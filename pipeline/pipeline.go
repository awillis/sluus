package pipeline

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor"
)

var (
	ErrInvalidProcessor = errors.New("invalid processor")
	ErrNoSource         = errors.New("missing source processor")
	ErrNoReject         = errors.New("missing reject sink processor")
	ErrNoAccept         = errors.New("missing accept sink processor")
)

type (
	Component struct {
		next  *Component
		pipe  *Pipe
		Value processor.Interface
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
		len       int
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
	pipe.logger = core.SetupLogger(core.LogConfig("pipeline", pipe.ID()))
	return pipe
}

func (p *Pipe) ID() string {
	return p.Id
}

func (p *Pipe) Logger() *zap.SugaredLogger {
	return p.logger
}

func (p *Pipe) Len() int {
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

func (p *Pipe) Run() {
	for n := &p.root; n.Next() != nil; n = n.Next() {
		n.Value.Run()
	}
}

func (p *Pipe) Attach(component *Component) {

	reject := processor.Reject(p.Reject().Value.Flume().Output())
	accept := processor.Accept(p.Accept().Value.Flume().Output())

	for n := &p.root; n != component; n = n.Next() {
		if n.Next() == nil {
			if component.Value.Type() != plugin.SINK {
				sluus := NewSluus()
				sluus.SetLogger(p.logger)

				if err := processor.Configure(sluus.Flume(), reject, accept); err != nil {
					p.Logger().Error(err)
				}

				tail := new(Component)
				tail.Value = sluus
				component.next = tail
				p.len++
			}

			component.Value.SetLogger(p.logger)
			if err := processor.Configure(component.Value.Flume(), reject, accept); err != nil {
				p.Logger().Error(err)
			}

			n.next = component
			n.pipe = p
		}

		if n.Next() != nil {
			switch n.Value.(type) {
			case *Sluus:
				n.Value.(*Sluus).SetReceiver(n.Next().Value)
			}
		}
	}
}

func (p *Pipe) AddSource(proc processor.Interface) (err error) {

	if proc.Type() != plugin.SOURCE {
		return ErrInvalidProcessor
	}

	src := new(Component)
	src.Value = proc
	p.Attach(src)
	p.hasSource = true
	p.len++
	return err
}

func (p *Pipe) AddConduit(proc processor.Interface) (err error) {

	if proc.Type() != plugin.CONDUIT {
		return ErrInvalidProcessor
	}

	if !p.hasSource {
		return ErrNoSource
	}

	conduit := new(Component)
	conduit.Value = proc
	p.Attach(conduit)
	p.len++
	return err
}

func (p *Pipe) AddReject(reject processor.Interface) (err error) {

	if reject.Type() != plugin.SINK {
		return ErrInvalidProcessor
	}

	if !p.hasSource {
		return ErrNoSource
	}

	sink := new(Component)
	sink.Value = reject
	p.Attach(sink)
	p.hasReject = true
	p.len++
	return err
}

func (p *Pipe) AddAccept(accept processor.Interface) (err error) {
	if accept.Type() != plugin.SINK {
		return ErrInvalidProcessor
	}

	if !p.hasSource {
		return ErrNoSource
	}

	if !p.hasReject {
		return ErrNoReject
	}

	sink := new(Component)
	sink.Value = accept
	p.Attach(sink)
	p.hasAccept = true
	p.len++
	return err
}
