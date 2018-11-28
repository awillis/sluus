package conduit

import "github.com/google/uuid"

type Conduit struct {
	id uuid.UUID
}

func NewConduit() Conduit {
	return Conduit{
		id: uuid.New(),
	}
}

func (s *Conduit) ID() uuid.UUID {
	return s.id
}
func (s *Conduit) Process() {

}
