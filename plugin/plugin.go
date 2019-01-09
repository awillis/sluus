package plugin

import (
	"fmt"
	"github.com/pkg/errors"
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
	Version() string
	Initialize() error
	Execute() error
	Shutdown() error
}

type Base struct {
	Interface
	Id       string
	PlugName string
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

func (b *Base) Version() string {
	return fmt.Sprintf("%d.%d.%d", b.Major, b.Minor, b.Patch)
}
