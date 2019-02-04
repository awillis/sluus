package processor

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/dgraph-io/badger/options"

	"github.com/awillis/sluus/message"
)

// design influenced by http://www.drdobbs.com/parallel/lock-free-queues/208801974

type Queue struct {
	sync.RWMutex
	opts     badger.Options
	db       *badger.DB
	readHead []byte
}

func NewQueue(dbPath string) (queue *Queue) {
	queue = new(Queue)
	queue.opts = badger.DefaultOptions

	// both keys and values can reside together
	queue.opts.Dir = dbPath
	queue.opts.ValueDir = dbPath
	// values are held in queue temporarily
	queue.opts.SyncWrites = false
	// the default value (mmap) assumes SSD
	queue.opts.ValueLogLoadingMode = options.FileIO
	// clear readHead
	queue.ResetReadHead()
	return
}

func (q *Queue) Initialize() (err error) {
	q.db, err = badger.Open(q.opts)
	return
}

func (q *Queue) Size() int64 {
	size, _ := q.db.Size()
	return size
}

func (q *Queue) ResetReadHead() {
	q.readHead = nil
	q.readHead = make([]byte, 0, 8)
}

func (q *Queue) Put(msgs []*message.Message) (err error) {

	err = q.db.Update(func(txn *badger.Txn) (e error) {

		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()

		key := new(bytes.Buffer)

		for _, msg := range msgs {
			binary.LittleEndian.PutUint64(key.Bytes(), uint64(time.Now().UnixNano()+msg.Received.GetSeconds()))
			e = txn.Set(key.Bytes(), []byte(msg.String()))
			key.Reset()
		}

		if q.Size() > 0 && len(q.readHead) > 0 {
			// if there is data present in the db and the read readHead is set
			// remove data from the beginning up to the read readHead
			for iter.Rewind(); iter.Valid(); iter.Next() {
				key := iter.Item().Key()
				if bytes.Equal(key, q.readHead) {
					break
				} else {
					e = txn.Delete(key)
				}
			}
		}
		return
	})
	return
}

func (q *Queue) Get(count int) (msgs []*message.Message, err error) {

	if q.Size() == 0 {
		return // no data, no error
	}

	err = q.db.View(func(txn *badger.Txn) (e error) {

		opts := badger.IteratorOptions{
			PrefetchValues: true,
			PrefetchSize:   count,
		}

		if opts.PrefetchSize > 128 {
			opts.PrefetchSize = 128
		}

		iter := txn.NewIterator(opts)
		defer iter.Close()

		// start at head if available, or at absolute start
		if len(q.readHead) > 0 {
			iter.Seek(q.readHead)
		} else {
			iter.Rewind()
		}

		// collect messages
		for i := count; iter.Valid() && i < count; i++ {

			var content []byte
			item := iter.Item()

			value, err := item.Value()
			if err != nil {
				e = err
			}

			copy(content, value)

			msg, err := message.WithContent(json.RawMessage(content))
			if err != nil {
				e = err
			}

			msgs = append(msgs, msg)
			iter.Next()
		}

		// if there are more records to be read, copy the key to seed the next read
		// otherwise clear the readHead so that the next read can start at the beginning
		if iter.Valid() {
			item := iter.Item()
			copy(q.readHead, item.Key())
		} else {
			q.ResetReadHead()
		}

		return
	})
	return
}

func (q *Queue) Shutdown() (err error) {
	return q.db.Close()
}
