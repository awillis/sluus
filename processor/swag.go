package processor

import (
	"github.com/awillis/sluus/core"
	"time"
)

type SWAG interface {
	Insert(val interface{}, tm time.Time)
	Evict(tm time.Time)
	Query(start time.Time, end time.Time, arg interface{})
}

type WBatch struct {
	core.Batch
}

func (wb *WBatch) Lift(msg core.Message) (*core.Message, int) {
	m := core.NewMessage(core.Message_NORMAL)
	return m, 1
}

func (wb *WBatch) Combine() {

}

func (wb *WBatch) Lower() {

}
