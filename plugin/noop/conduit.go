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

	rCount := uint64(float64(c.opts.RejectPercentage) * 0.01 * float64(input.Count()))
	aCount := uint64(float64(c.opts.AcceptPercentage) * 0.01 * float64(input.Count()))

	reject = message.NewBatch(rCount)
	accept = message.NewBatch(aCount)

	//for msg := range input.Iter() {
	//	switch {
	//	case reject.Count() <= rCount:
	//		if e := reject.Add(msg); e == message.ErrBatchFull {
	//			input.Cancel()
	//			continue
	//		}
	//	case accept.Count() <= aCount:
	//		if e := accept.Add(msg); e != nil {
	//			input.Cancel()
	//			continue
	//		}
	//	}
	//}
	c.Logger().Infof("conduit return input: %d, reject: %d, accept: %d", input.Count(), reject.Count(), accept.Count())
	return input, reject, accept
}

func (c *Conduit) Shutdown() (err error) {
	return
}
