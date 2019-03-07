package plugin

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/awillis/sluus/message"
)

type Type uint8

const (
	CONDUIT Type = iota
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
		Start(ctx context.Context)
		Produce() (ch <-chan *message.Batch)
		Shutdown() (err error)
	}

	Processor interface {
		Start(ctx context.Context)
		Process(input *message.Batch) (output *message.Batch)
		Shutdown() (err error)
	}

	Consumer interface {
		Start(ctx context.Context)
		Consume(batch *message.Batch)
		Shutdown() (err error)
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
