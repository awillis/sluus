package noop

import (
	"context"
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

func (c *Conduit) Execute() (err error) {
	return
}

func (c *Conduit) Shutdown() (err error) {
	return
}
