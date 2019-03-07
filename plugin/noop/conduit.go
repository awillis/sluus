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

func (c *Conduit) Process(batch *message.Batch) (pbatch *message.Batch) {

	rCount := uint64(float64(c.opts.RejectPercentage) * 0.01 * float64(batch.Count()))
	aCount := uint64(float64(c.opts.AcceptPercentage) * 0.01 * float64(batch.Count()))

	pbatch = message.NewBatch(batch.Count())

	for msg := range batch.Iter() {

		if rCount > 0 {
			msg.Direction = message.Route_REJECT
			pbatch.Add(msg)
			rCount--
			continue
		}

		if aCount > 0 {
			msg.Direction = message.Route_ACCEPT
			pbatch.Add(msg)
			aCount--
			continue
		}

		pbatch.Add(msg)
	}

	c.Logger().Infof("conduit return batch total: %d, output: %d, reject: %d, accept: %d",
		pbatch.Count(),
		pbatch.PassCount(),
		pbatch.AcceptCount(),
		pbatch.RejectCount())
	return pbatch
}

func (c *Conduit) Shutdown() (err error) {
	return
}
