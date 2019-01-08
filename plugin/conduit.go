package plugin

import (
	"runtime"
)

type Conduit struct {
	Plugin
}

func (c *Conduit) Run() {

}

func (c *Conduit) Execute() error {
	runtime.LockOSThread()
	var err error
	return err
}

func (c *Conduit) Shutdown() error {
	var err error
	return err
}
