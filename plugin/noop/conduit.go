package noop

import (
	"context"
	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

var _ plugin.Interface = new(Conduit)
var _ plugin.Processor = new(Conduit)

type Conduit struct {
	plugin.Base
	opts *options
}

func (c *Conduit) Options() interface{} {
	return c.opts
}

func (c *Conduit) Initialize() (err error) {
	plugin.Validate(c.opts,
		c.opts.defaultMessagePerBatch(),
		c.opts.defaultRejectPercentage(),
		c.opts.defaultAcceptPercentage(),
	)
	c.opts.logCurrentConfig(c.Logger())
	return
}

func (c *Conduit) Start(ctx context.Context) {
	return
}

func (c *Conduit) Process(input *message.Batch) (output, reject, accept *message.Batch) {

	rCount := uint64(float64(c.opts.RejectPercentage/100) * float64(input.Count()))
	aCount := uint64(float64(c.opts.AcceptPercentage/100) * float64(input.Count()))

	reject = message.NewBatch(rCount)
	accept = message.NewBatch(aCount)

	for msg := range input.Iter() {
		switch {
		case reject.Count() <= rCount:
			if e := reject.Add(msg); e != nil {
				c.Logger().Error(e)
			}
		case accept.Count() <= aCount:
			if e := accept.Add(msg); e != nil {
				c.Logger().Error(e)
			}
		default:
			input.Cancel()
		}
	}
	c.Logger().Infof("sending output %d records", input.Count())
	return input, reject, accept
}

func (c *Conduit) Shutdown() (err error) {
	return
}
