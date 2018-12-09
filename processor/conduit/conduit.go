package conduit

import (
	"sluus/plugin"
)

type Conduit struct {
	plugin.PlugBase
}

type Conveyor interface {
	Convey()
}

func (c *Conduit) Convey() {

}
