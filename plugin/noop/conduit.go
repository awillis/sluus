package noop

import (
	"context"
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

var _ plugin.Interface = new(Conduit)
var _ plugin.Processor = new(Conduit)

type Conduit struct {
	plugin.Base
	opts *options
}

func (c *Conduit) Options() interface{} {
	return c.opts
}

func (c *Conduit) Initialize() (err error) {
	return
}

func (c *Conduit) Start(ctx context.Context) {
	return
}

func (c *Conduit) Process(input *message.Batch) (output, reject, accept *message.Batch, err error) {
	output = new(message.Batch)
	reject = new(message.Batch)
	accept = new(message.Batch)
	return
}

func (c *Conduit) Shutdown() (err error) {
	return
}
