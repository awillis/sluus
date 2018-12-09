package plugin

import (
	"fmt"
	"plugin"
)

type Plugin interface {
	Load(name string) bool
	Initialize() error
	Execute() error
	Shutdown() error
}

type PlugBase struct {
	Name       string
	Version    string
	initialize func() error
	execute    func() error
	shutdown   func() error
}

func (p *PlugBase) Load(name string) bool {

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

	p.Name = symName.(string)

	symVer, err := plug.Lookup("version")

	if err != nil {
		fmt.Println(err)
		return false
	}

	p.Version = symVer.(string)

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

func (p *PlugBase) Initialize() error {
	return p.initialize()
}

func (p *PlugBase) Execute() error {
	return p.execute()
}

func (p *PlugBase) Shutdown() error {
	return p.shutdown()
}
