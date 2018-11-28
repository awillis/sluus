package plugin

import "github.com/google/uuid"

type Processor interface {
	ID() uuid.UUID
	Process()
}
