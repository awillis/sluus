package plugin

import (
	"fmt"
	"os"
	"plugin"
	"strings"

	"github.com/awillis/sluus/core"
)

func (p *Plugin) Load(plugname string) bool {

	plugfile := strings.Join([]string{core.PLUGDIR, plugname + ".so"}, string(os.PathSeparator))
	plug, err := plugin.Open(plugfile)

	if err != nil {
		_ = fmt.Errorf("error loading plugin: %v", err)
		return false
	}

	symName, err := plug.Lookup("name")

	if err != nil {
		fmt.Println(err)
		return false
	}

	p.name = symName.(string)

	symVer, err := plug.Lookup("version")

	if err != nil {
		fmt.Println(err)
		return false
	}

	p.version = symVer.(string)

	symInit, err := plug.Lookup("initialize")

	if err != nil {
		fmt.Println(err)
		return false
	}

	p.initialize = symInit.(func() error)

	symExec, err := plug.Lookup("execute")

	if err != nil {
		fmt.Println(err)
		return false
	}

	p.execute = symExec.(func() error)

	symShut, err := plug.Lookup("shutdown")

	if err != nil {
		fmt.Println(err)
		return false
	}

	p.shutdown = symShut.(func() error)

	return true
}
