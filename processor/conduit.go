package processor

import (
	"github.com/awillis/sluus/plugin"
)

type Conduit struct {
	plugin.PlugBase
}

type Conveyor interface {
	Convey()
}

func (c *Conduit) Convey() {

}
