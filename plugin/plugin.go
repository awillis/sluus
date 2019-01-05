package plugin

import (
	"fmt"
	"plugin"
)

const (
	PROCESSOR ComponentType = iota
	MESSAGE
)

type ComponentType int

type Component interface {
	Name() string
	Version() string
	Load(name string) bool
	Initialize() error
}

type Plugin struct {
	Component
	name       string
	version    string
	initialize func() error
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

	return true
}
