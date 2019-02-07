package tcp

import "github.com/awillis/sluus/plugin"

func init() {
	plugin.Registry.Register("tcp", New)
}
