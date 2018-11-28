package sink

import "github.com/google/uuid"

type Sink struct {
	id uuid.UUID
}

func NewSink() Sink {
	return Sink{
		id: uuid.New(),
	}
}

func (s *Sink) ID() uuid.UUID {
	return s.id
}
func (s *Sink) Process() {

}
