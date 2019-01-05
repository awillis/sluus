package processor

import (
	"github.com/awillis/sluus/plugin"
)

type Conduit struct {
	plugin.Plugin
}

func (c *Conduit) Run() {

}

func (c *Conduit) Execute() error {
	var err error
	return err
}

func (c *Conduit) Shutdown() error {
	var err error
	return err
}
