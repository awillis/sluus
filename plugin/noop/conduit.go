package noop

import (
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

func (c *Conduit) Initialize() (err error) {
	return
}

func (c *Conduit) Execute(input <-chan message.Batch, accept chan<- message.Batch, reject chan<- message.Batch) (err error) {
	return
}

func (c *Conduit) Shutdown() (err error) {
	return
}
