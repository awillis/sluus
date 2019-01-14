package plugin

import (
	"errors"
	"fmt"
)

type Type uint8

const (
	MESSAGE Type = iota
	CONDUIT
	SOURCE
	SINK
)

var ErrUnimplemented = errors.New("unimplemented plugin")

type Interface interface {
	ID() string
	Name() string
	Type() Type
	Version() string
}

type Processor interface {
	Interface
	Initialize() error
	Execute() error
	Shutdown() error
}

type Base struct {
	Id       string
	PlugName string
	PlugType Type
	Major    uint8
	Minor    uint8
	Patch    uint8
}

func (b *Base) ID() string {
	return b.Id
}

func (b *Base) Name() string {
	return b.PlugName
}

func (b *Base) Type() Type {
	return b.PlugType
}

func (b *Base) Version() string {
	return fmt.Sprintf("%d.%d.%d", b.Major, b.Minor, b.Patch)
}
