package processor

import (
	"github.com/awillis/sluus/plugin"
)

type Source struct {
	plugin.Plugin
}

func (s *Source) Run() {

}

func (s *Source) Execute() error {
	var err error
	return err
}

func (s *Source) Shutdown() error {
	var err error
	return err
}

// TODO: Create a heap.Interface for source processors to use in sorting generated messages by priority
// TODO: receive messages in a constant thread into the priority heap
// TODO: create go thread that constantly sorts the priority heap and creates batches
