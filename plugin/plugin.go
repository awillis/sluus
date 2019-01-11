package plugin

import (
	"errors"
)

type Type uint8

const (
	CONDUIT Type = iota
	SOURCE
	SINK
)

var ErrUnimplemented = errors.New("unimplemented plugin")

type Interface interface {
	ID() string
	Name() string
	Type() Type
	Version() string
	Initialize() error
	Execute() error
	Shutdown() error
}

type Base struct {
	Interface
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
	return string(b.Major) + "." + string(b.Minor) + "." + string(b.Patch)
}
