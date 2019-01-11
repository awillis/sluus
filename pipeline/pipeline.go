package pipeline

import (
	"errors"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/awillis/sluus/core"
	"github.com/awillis/sluus/plugin"
	"github.com/awillis/sluus/processor"
)

var ErrInvalidProcessor = errors.New("invalid processor")
var ErrMissSourceSink = errors.New("missing source or sink processor")

type Component struct {
	next, prev *Component
	pipe       *Pipe
	Value      processor.Interface
}

func (c *Component) Next() *Component {
	if p := c.next; c.pipe != nil && p != &c.pipe.root {
		return p
	}
	return nil
}

func (c *Component) Prev() *Component {
	if p := c.prev; c.pipe != nil && p != &c.pipe.root {
		return p
	}
	return nil
}

type Pipe struct {
	Id        string
	logger    *zap.SugaredLogger
	hasSource bool
	hasSink   bool
	root      Component
	len       int
}

func NewPipeline() *Pipe {

	pipe := new(Pipe)
	pipe.Id = uuid.New().String()
	pipe.root.next = &pipe.root
	pipe.root.prev = &pipe.root
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

func (p *Pipe) Sink() *Component {
	if p.len == 0 {
		return nil
	}
	return p.root.prev
}

func (p *Pipe) SetSource(proc processor.Interface) error {

	if proc.Type() != plugin.SOURCE {
		return ErrInvalidProcessor
	}

	src := new(Component)
	proc.SetLogger(p.logger)
	src.Value = proc
	src.pipe = p
	src.prev = &p.root
	p.root.next = src
	p.hasSource = true
	p.len++
	return nil
}

func (p *Pipe) SetSink(proc processor.Interface) error {
	if proc.Type() != plugin.SINK {
		return ErrInvalidProcessor
	}

	sink := new(Component)
	proc.SetLogger(p.logger)
	sink.Value = proc
	sink.pipe = p
	sink.next = &p.root
	p.root.prev = sink
	p.hasSink = true
	p.len++
	return nil
}

func (p *Pipe) AddConduit(proc processor.Interface) error {

	if proc.Type() != plugin.CONDUIT {
		return ErrInvalidProcessor
	}

	if !p.hasSource || !p.hasSink {
		return ErrMissSourceSink
	}

	conduit := new(Component)
	proc.SetLogger(p.logger)
	conduit.Value = proc
	conduit.pipe = p
	prev := p.Sink().prev
	p.Sink().prev = conduit
	conduit.prev = prev
	conduit.next = p.Sink()
	prev.next = conduit

	p.len++
	return nil
}
