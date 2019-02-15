package processor

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"go.uber.org/zap"
	"os"
	"sync"
	"time"

	"github.com/awillis/sluus/message"
	"github.com/dgraph-io/badger"
)

// design influenced by http://www.drdobbs.com/parallel/lock-free-queues/208801974

const (
	INPUT uint64 = iota
	OUTPUT
	REJECT
	ACCEPT
)

type (
	Queue struct {
		batchSize uint64
		opts      badger.Options
		db        *badger.DB
		head      readHead
		cancel    context.CancelFunc
		logger    *zap.SugaredLogger
	}

	readHead struct {
		sync.RWMutex
		m map[uint64][]byte
	}

	QueueOpt func(*Queue) error
)

func NewQueue() (queue *Queue) {
	queue = new(Queue)
	queue.opts = badger.DefaultOptions
	queue.opts.SyncWrites = false
	queue.head = readHead{
		m: make(map[uint64][]byte),
	}

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
	if e := os.MkdirAll(q.opts.Dir, 0755); e != nil {
		return e
	}
	q.db, err = badger.Open(q.opts)
	return
}

func (q *Queue) Logger() *zap.SugaredLogger {
	return q.logger.With("queue")
}

func (q *Queue) Size() int64 {
	size, _ := q.db.Size()
	return size
}

func u64ToBytes(i uint64) (b []byte) {
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return
}

func (q *Queue) Put(prefix uint64, batch *message.Batch) {

	err := q.db.Update(func(txn *badger.Txn) (e error) {

		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()

		hash := md5.New()
		prefixKey := u64ToBytes(prefix)
		timeKey := u64ToBytes(uint64(time.Now().UnixNano()))
		key := make([]byte, len(prefixKey)+len(timeKey)+md5.Size)

		for msg := range batch.Iter() {

			contentKey := hash.Sum([]byte(msg.String()))
			key = append(key, prefixKey...)
			key = append(key, timeKey...)
			key = append(key, contentKey...)

			e = txn.Set(key, []byte(msg.String()))
			key = key[:0]
			hash.Reset()
		}

		if q.Size() > 0 && len(q.head.Get(prefix)) > 0 {
			// if there is data present in the db and the read head is set
			// remove data from the beginning up to the read head
			for iter.Rewind(); iter.ValidForPrefix(prefixKey); iter.Next() {
				key := iter.Item().Key()
				if bytes.Equal(key, q.head.Get(prefix)) {
					break
				} else {
					e = txn.Delete(key)
				}
			}
		}
		return
	})

	if err != nil {
		q.Logger().Error(err)
	}

	return
}

func (q *Queue) Get(prefix, size uint64) <-chan *message.Batch {
	iter := make(chan *message.Batch)
	ctx := context.Background()
	ctx, cancel := context.WithCancel(context.Background())
	q.cancel = cancel

	go q.query(ctx, iter, prefix, size)
	return iter
}

func (q *Queue) query(ctx context.Context, iter chan *message.Batch, prefix, size uint64) {

	defer close(iter)
	prefixKey := u64ToBytes(prefix)
	batch := message.NewBatch(q.batchSize)

	if size == 0 || size > q.batchSize {
		size = q.batchSize
	}

	err := q.db.View(func(txn *badger.Txn) (e error) {

		opts := badger.IteratorOptions{
			PrefetchValues: true,
			PrefetchSize:   int(size),
		}

		if opts.PrefetchSize > 128 {
			opts.PrefetchSize = 128
		}

		iter := txn.NewIterator(opts)
		defer iter.Close()

		// start at head if available, or at absolute start
		if len(q.head.Get(prefix)) > 0 {
			iter.Seek(q.head.Get(prefix))
		} else {
			iter.Rewind()
		}

		// collect messages

	cancel:
		for i := size; iter.ValidForPrefix(prefixKey) && i < size; i++ {

			select {
			case <-ctx.Done():
				break cancel
			default:
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

				if err := batch.Add(msg); err != nil {
					break
				}
				iter.Next()
			}
		}

		// if there are more records to be read, copy the key to seed the next read
		// otherwise clear the head so that the next read can start at the beginning
		if iter.ValidForPrefix(prefixKey) {
			item := iter.Item()
			q.head.Set(prefix, item.Key())
		} else {
			q.head.Reset(prefix)
		}

		return
	})

	if err != nil {
		q.Logger().Error(err)
	}
	iter <- batch
}

func (q *Queue) Cancel() {
	q.cancel()
}

func (q *Queue) shutdown() (err error) {
	return q.db.Close()
}

func (rh *readHead) Get(prefix uint64) []byte {
	rh.Lock()
	defer rh.Unlock()
	return rh.m[prefix]
}

func (rh *readHead) Set(prefix uint64, key []byte) {
	rh.Lock()
	defer rh.Unlock()
	copy(rh.m[prefix], key)
}

func (rh *readHead) Reset(prefix uint64) {
	rh.Lock()
	defer rh.Unlock()
	rh.m[prefix] = nil
	rh.m[prefix] = make([]byte, 32)
}
