package processor

import (
	"github.com/awillis/sluus/message"
	"time"
)

type SWAG interface {
	Insert(val interface{}, tm time.Time)
	Evict(tm time.Time)
	Query(start time.Time, end time.Time, arg interface{})
}

type WBatch struct {
	message.Batch
}

func (wb *WBatch) Lift(msg message.Message) (*message.Message, int) {
	m := message.NewMessage(message.Message_NORMAL)
	return m, 1
}

func (wb *WBatch) Combine() {

}

func (wb *WBatch) Lower() {

}
