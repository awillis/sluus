package noop

import (
	"context"
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

type Conduit struct {
	plugin.Base
	opts options
}

func (c *Conduit) Options() interface{} {
	return &c.opts
}

func (c *Conduit) Initialize(ctx context.Context) (err error) {
	return
}

func (c *Conduit) Process(message.Batch) (batch message.Batch, err error) {
	return batch
}

func (c *Conduit) Shutdown() (err error) {
	return
}
