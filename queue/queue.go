package queue

import (
	"encoding/json"
	"github.com/awillis/sluus/message"
	"github.com/dgraph-io/badger"
	"sync"
)

// design influenced by http://www.drdobbs.com/parallel/lock-free-queues/208801974

type Queue struct {
	sync.RWMutex
	opts       badger.Options
	db         *badger.DB
	head, tail []byte
}

func New(dbPath string) (queue *Queue) {
	queue = new(Queue)
	queue.opts = badger.DefaultOptions
	queue.opts.Dir = dbPath
	queue.opts.ValueDir = dbPath
	return
}

func (q *Queue) Initialize() (err error) {
	q.db, err = badger.Open(q.opts)
	// put one record in the db
	return
}

func (q *Queue) Produce(msg *message.Message) (err error) {

	err = q.db.Update(func(txn *badger.Txn) (e error) {

		opts := badger.IteratorOptions{
			PrefetchValues: true,
			PrefetchSize:   64,
		}

		iter := txn.NewIterator(opts)
		iter.Item().KeyCopy(q.head)

		e = txn.Set([]byte(msg.Id), []byte(msg.String()))
		q.tail = []byte(msg.Id)

		for iter.Rewind(); iter.Valid(); iter.Next() {
			// delete from the very start of the transaction to the head key
			key := iter.Item().Key()
			if string(key) != string(q.head) {
				e = txn.Delete(key)
			} else {
				break
			}
		}
		return
	})
	return
}

func (q *Queue) Consume() (msg *message.Message, err error) {

	var content []byte

	err = q.db.View(func(txn *badger.Txn) (e error) {

		opts := badger.IteratorOptions{
			PrefetchValues: true,
			PrefetchSize:   1,
		}

		iter := txn.NewIterator(opts)
		iter.Seek(q.head)
		iter.Next()

		if next := iter.Item(); string(next.Key()) != string(q.tail) {
			value, e := next.Value()
			if e != nil {
				return e
			}
			q.head = next.Key()
			copy(content, value)
		}
		return
	})

	return message.WithContent(json.RawMessage(content))
}

func (q *Queue) Shutdown() (err error) {
	return q.db.Close()
}
