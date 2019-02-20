package noop

import (
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
	"github.com/google/uuid"
	"sync"
)

const (
	NAME  string = "noop"
	MAJOR uint8  = 0
	MINOR uint8  = 0
	PATCH uint8  = 1
)

type options struct {
	BatchInterval   uint64 `toml:"batch_interval"`
	MessagePerBatch uint64 `toml:"message_per_batch"`
}

func New(pluginType plugin.Type) (plug plugin.Interface, err error) {

	switch pluginType {
	case plugin.SOURCE:
		return &Source{
			opts:   new(options),
			wg:     new(sync.WaitGroup),
			output: make(chan *message.Batch),
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: NAME,
				PlugType: pluginType,
				Major:    MAJOR,
				Minor:    MINOR,
				Patch:    PATCH,
			},
		}, err
	case plugin.CONDUIT:
		return &Conduit{
			opts: new(options),
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: NAME,
				PlugType: pluginType,
				Major:    MAJOR,
				Minor:    MINOR,
				Patch:    PATCH,
			},
		}, err
	case plugin.SINK:
		return &Sink{
			opts: new(options),
			Base: plugin.Base{
				Id:       uuid.New().String(),
				PlugName: NAME,
				PlugType: pluginType,
				Major:    MAJOR,
				Minor:    MINOR,
				Patch:    PATCH,
			},
		}, err
	default:
		return plug, plugin.ErrUnimplemented
	}
}

func (o *options) validMessagePerBatch() plugin.Default {
	return func(def plugin.Option) {
		if o.MessagePerBatch == 0 {
			o.MessagePerBatch = 5
		}
	}
}

func (o *options) validBatchInterval() plugin.Default {
	return func(def plugin.Option) {
		if o.BatchInterval == 0 {
			o.BatchInterval = 5
		}
	}
}
