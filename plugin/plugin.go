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
		SetLogger(*zap.SugaredLogger)
		Logger() *zap.SugaredLogger
	}

	Producer interface {
		Produce() chan *message.Batch
		Shutdown() error
	}

	Processor interface {
		Process(*message.Batch) (output, reject, accept *message.Batch, err error)
		Shutdown() error
	}

	Consumer interface {
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
		logger   *zap.SugaredLogger
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
	b.logger = logger
}

func (b *Base) Logger() *zap.SugaredLogger {
	return b.logger.With("type", TypeName(b.Type()), "plugin_id", b.ID(), "plugin", b.Name())
}

func TypeName(t Type) (s string) {
	switch t {
	case SOURCE:
		s = "source"
	case SINK:
		s = "sink"
	case CONDUIT:
		s = "conduit"
	}
	return
}
