package processor

import (
	"github.com/awillis/sluus/plugin"
)

type Conduit struct {
	conduit
	plugin.Plugin
}

type conduit interface {
	Convey()
}

func (c *Conduit) Convey() {

}
