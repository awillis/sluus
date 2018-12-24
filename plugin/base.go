package plugin

import (
	"fmt"
	"plugin"
)

type Component interface {
	Name() string
	Version() string
	Load(name string) bool
	Initialize() error
	Execute() error
	Shutdown() error
}

type Plugin struct {
	Component
	name       string
	version    string
	initialize func() error
	execute    func() error
	shutdown   func() error
}

func (p *Plugin) Load(name string) bool {

	plug, err := plugin.Open(name)

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

func (p Plugin) Name() string {
	return p.name
}

func (p Plugin) Version() string {
	return p.version
}

func (p *Plugin) Initialize() error {
	return p.initialize()
}

func (p *Plugin) Execute() error {
	return p.execute()
}

func (p *Plugin) Shutdown() error {
	return p.shutdown()
}
