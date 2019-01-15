package noop

import "github.com/awillis/sluus/plugin"

type Conduit struct {
	plugin.Base
}

func (c *Conduit) Initialize() (err error) {
	return
}

func (c *Conduit) Execute() (err error) {
	return
}

func (c *Conduit) Shutdown() (err error) {
	return
}
