package source

import "github.com/google/uuid"

type Source struct {
	id uuid.UUID
}

func NewSource() Source {
	return Source{
		id: uuid.New(),
	}
}

func (s *Source) ID() uuid.UUID {
	return s.id
}
func (s *Source) Process() {

}
