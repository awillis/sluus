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

func (c *Conduit) Process(input *message.Batch) (output *message.Batch) {

	rCount := uint64(float64(c.opts.RejectPercentage) * 0.01 * float64(input.Count()))
	aCount := uint64(float64(c.opts.AcceptPercentage) * 0.01 * float64(input.Count()))

	output = message.NewBatch(input.Count())

	for msg := range input.Iter() {

		if rCount > 0 {
			msg.Direction = message.Route_REJECT
			output.Add(msg)
			rCount--
			continue
		}

		if aCount > 0 {
			msg.Direction = message.Route_ACCEPT
			output.Add(msg)
			aCount--
			continue
		}

		output.Add(msg)
	}

	c.Logger().Infof("conduit return input total: %d, output: %d, reject: %d, accept: %d",
		output.Count(),
		output.PassCount(),
		output.AcceptCount(),
		output.RejectCount())
	return output
}

func (c *Conduit) Shutdown() (err error) {
	return
}
