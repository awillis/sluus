package plugin

import (
	"fmt"
	"go.uber.org/zap"

	"github.com/awillis/sluus/message"
	"github.com/pkg/errors"
)

type Type uint8

const (
	MESSAGE Type = iota
	CONDUIT
	SOURCE
	SINK
)

var ErrUnimplemented = errors.New("unimplemented plugin")

type (
	Interface interface {
		ID() string
		Name() string
		Type() Type
		Version() string
		Options() interface{}
		Initialize() error
	}

	Producer interface {
		SetLogger(*zap.SugaredLogger)
		Produce() chan *message.Batch
		Shutdown() error
	}

	Processor interface {
		SetLogger(*zap.SugaredLogger)
		Process(*message.Batch) (output, reject, accept *message.Batch, err error)
		Shutdown() error
	}

	Consumer interface {
		SetLogger(*zap.SugaredLogger)
		Consume() chan *message.Batch
		Shutdown() error
	}

	Base struct {
		Id       string
		PlugName string
		PlugType Type
		Major    uint8
		Minor    uint8
		Patch    uint8
		Logger   *zap.SugaredLogger
	}
)

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

func (b *Base) SetLogger(logger *zap.SugaredLogger) {
	b.Logger = logger
}
