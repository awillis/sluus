package processor

import (
	"github.com/golang-collections/go-datastructures/queue"
	"kapilary/core"
)

type Base struct {
	input chan core.Event
	output chan core.Event
	queue queue.PriorityQueue
	category string
}

func (b Base) Input() chan core.Event {
	return b.input
}

func (b Base) Output() chan core.Event {
	return b.output
}

func (b Base) Queue() *queue.PriorityQueue {
	return &b.queue
}

func (b Base) Category() string {
	return b.category
}