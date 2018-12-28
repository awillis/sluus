package plugin

import (
	"fmt"
	"plugin"
)

type ProcessorComponent interface {
	Name() string
	Version() string
	Load(name string) bool
	Initialize() error
	Execute() error
	Shutdown() error
	Run()
}

type ProcessorPlugin struct {
	ProcessorComponent
	name       string
	version    string
	initialize func() error
	execute    func() error
	shutdown   func() error
}

func (p *ProcessorPlugin) Load(name string) bool {

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

func (p ProcessorPlugin) Name() string {
	return p.name
}

func (p ProcessorPlugin) Version() string {
	return p.version
}

func (p *ProcessorPlugin) Initialize() error {
	return p.initialize()
}

func (p *ProcessorPlugin) Execute() error {
	return p.execute()
}

func (p *ProcessorPlugin) Shutdown() error {
	return p.shutdown()
}
