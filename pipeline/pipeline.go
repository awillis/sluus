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
	pipe.logger = core.SetupLogger(core.LogConfig("pipeline", pipe.ID()))
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

func (p *Pipe) Start() {
	for n := &p.root; n.Next() != nil; n = n.Next() {
		n.Value.Start()
	}
}

func (p *Pipe) Stop() {
	for n := &p.root; n.Next() != nil; n = n.Next() {
		n.Value.Stop()
	}
}

func (p *Pipe) Attach(component *Component) {

	for n := &p.root; n != component; n = n.Next() {
		if n.Next() == nil {
			component.Value.SetLogger(p.logger)
			if e := component.Value.Initialize(); e != nil {
				p.Logger().Errorw(e.Error(),
					"name", component.Value.Plugin().Name(),
					"id", component.Value.ID())
			}
			p.len++
			n.pipe = p
			n.next = component
		}
	}
}

func (p *Pipe) Configure() {
	reject := processor.Reject(p.Reject().Value.Sluus().Input())
	accept := processor.Accept(p.Accept().Value.Sluus().Input())

	for n := &p.root; n.Next() == nil; n = n.Next() {
		if err := processor.Configure(n.Value.Sluus(), reject, accept); err != nil {
			p.Logger().Error(err)
		}

		if n.Next() != p.Accept() {
			input := processor.Input(n.Value.Sluus().Output())
			if err := processor.Configure(n.Next().Value.Sluus(), input); err != nil {
				p.Logger().Error(err)
			}
		}
	}
}

func (p *Pipe) Add(proc processor.Interface) (err error) {

	component := new(Component)
	component.Value = proc

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
