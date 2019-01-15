package plugin

import "errors"

var ErrInvalidOption = errors.New("invalid option")

type Option func(Processor) error

func Configure(processor Processor, opts ...Option) (err error) {
	for _, o := range opts {
		err = o(processor)
		if err != nil {
			return
		}
	}
	return
}
