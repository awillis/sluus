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

func (wb *WBatch) Lift(msg message.Message) (mssg *message.Message, i int) {
	mssg = message.New(message.Message_NORMAL)
	return mssg, 1
}

func (wb *WBatch) Combine() {

}

func (wb *WBatch) Lower() {

}
