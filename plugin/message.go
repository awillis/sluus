package plugin

import (
	"fmt"
	"plugin"
)

type Message interface {
	Name() string
	Version() string
	Load(name string) bool
	Initialize() error
}

type BasicJSON struct {
	Message
	name       string
	version    string
	initialize func() error
}

func (m *BasicJSON) Load(name string) bool {

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

	m.name = symName.(string)

	symVer, err := plug.Lookup("version")

	if err != nil {
		fmt.Println(err)
		return false
	}

	m.version = symVer.(string)

	symInit, err := plug.Lookup("initialize")

	if err != nil {
		fmt.Println(err)
		return false
	}

	m.initialize = symInit.(func() error)

	return true
}

func (m BasicJSON) Name() string {
	return m.name
}

func (m BasicJSON) Version() string {
	return m.version
}

func (m *BasicJSON) Initialize() error {
	return m.initialize()
}
