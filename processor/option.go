package processor

// configuration options

func Configure(sluus *Sluus, opts ...Option) (err error) {
	for _, o := range opts {
		err = o(sluus)
		if err != nil {
			return
		}
	}
	return
}

func Input(input chan *message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.input = input
		return
	}
}

func Output(output chan *message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.output = output
		return
	}
}

func Reject(reject chan *message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.reject = reject
		return
	}
}

func Accept(accept chan *message.Batch) Option {
	return func(s *Sluus) (err error) {
		s.accept = accept
		return
	}
}
