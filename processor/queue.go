package processor

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/awillis/sluus/message"
	"github.com/dgraph-io/badger"
)

// design influenced by http://www.drdobbs.com/parallel/lock-free-queues/208801974

const (
	INPUT byte = 3 << iota
	OUTPUT
	REJECT
	ACCEPT
)

type (
	Queue struct {
		opts     badger.Options
		db       *badger.DB
		readHead map[byte][]byte
	}

	QueueOpt func(*Queue) error
)

func NewQueue() (queue *Queue) {
	queue = new(Queue)
	queue.opts = badger.DefaultOptions
	queue.opts.SyncWrites = false
	queue.readHead = make(map[byte][]byte)
	return
}

func (q *Queue) Configure(opts ...QueueOpt) (err error) {
	for _, o := range opts {
		err = o(q)
		if err != nil {
			return
		}
	}
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

func (q *Queue) resetHead(prefix byte) {
	q.readHead[prefix] = nil
	q.readHead[prefix] = make([]byte, 0, 8)
}

func (q *Queue) Put(prefix byte, batch *message.Batch) (err error) {

	err = q.db.Update(func(txn *badger.Txn) (e error) {

		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()

		key := new(bytes.Buffer)
		key.WriteByte(prefix)

		for msg := range batch.Iter() {
			binary.LittleEndian.PutUint64(key.Bytes(), uint64(time.Now().UnixNano()+msg.Received.GetSeconds()))
			e = txn.Set(key.Bytes(), []byte(msg.String()))
			key.Reset()
		}

		if q.Size() > 0 && len(q.readHead) > 0 {
			// if there is data present in the db and the read readHead is set
			// remove data from the beginning up to the read readHead
			for iter.Rewind(); iter.ValidForPrefix([]byte{prefix}); iter.Next() {
				key := iter.Item().Key()
				if bytes.Equal(key, q.readHead[prefix]) {
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

func (q *Queue) Get(prefix byte, batchSize uint) (batch *message.Batch, err error) {

	if q.Size() == 0 {
		return // no data, no error
	}

	err = q.db.View(func(txn *badger.Txn) (e error) {

		opts := badger.IteratorOptions{
			PrefetchValues: true,
			PrefetchSize:   int(batchSize),
		}

		if opts.PrefetchSize > 128 {
			opts.PrefetchSize = 128
		}

		iter := txn.NewIterator(opts)
		defer iter.Close()

		// start at head if available, or at absolute start
		if len(q.readHead[prefix]) > 0 {
			iter.Seek(q.readHead[prefix])
		} else {
			iter.Rewind()
		}

		batch := message.NewBatch(batchSize)

		// collect messages
		for i := batchSize; iter.ValidForPrefix([]byte{prefix}) && i < batchSize; i++ {

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

			_ = batch.Add(msg)
			iter.Next()
		}

		// if there are more records to be read, copy the key to seed the next read
		// otherwise clear the readHead so that the next read can start at the beginning
		if iter.ValidForPrefix([]byte{prefix}) {
			item := iter.Item()
			copy(q.readHead[prefix], item.Key())
		} else {
			q.resetHead(prefix)
		}

		return
	})
	return
}

func (q *Queue) shutdown() (err error) {
	return q.db.Close()
}
