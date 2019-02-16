package processor

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"github.com/awillis/sluus/plugin"
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
	queue struct {
		pType        plugin.Type
		batchSize    uint64
		opts         badger.Options
		db           *badger.DB
		head         head
		cancel       context.CancelFunc
		requestChan  map[uint64]chan uint64
		responseChan map[uint64]chan *message.Batch
		logger       *zap.SugaredLogger
	}

	head struct {
		sync.RWMutex
		m map[uint64][]byte
	}
)

func newQueue(pType plugin.Type) (q *queue) {
	q = new(queue)
	q.pType = pType
	q.opts = badger.DefaultOptions
	q.opts.SyncWrites = false
	q.requestChan = make(map[uint64]chan uint64)
	q.responseChan = make(map[uint64]chan *message.Batch)
	q.head = head{
		m: make(map[uint64][]byte),
	}
	return
}

func (q *queue) Initialize() (err error) {

	if e := os.MkdirAll(q.opts.Dir, 0755); e != nil {
		return e
	}
	q.db, err = badger.Open(q.opts)
	return
}

func (q *queue) Start(ctx context.Context) {
	go q.query(ctx, INPUT)
	go q.query(ctx, OUTPUT)
	go q.query(ctx, ACCEPT)
	go q.query(ctx, REJECT)
}

func (q *queue) Logger() *zap.SugaredLogger {
	return q.logger.With("queue")
}

func u64ToBytes(i uint64) (b []byte) {
	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, i)
	return
}

func (q *queue) Put(prefix uint64, batch *message.Batch) {

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

		if len(q.head.Get(prefix)) > 0 {
			// if the read head is set, remove data
			// from the beginning up to the read head
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

func (q *queue) Get(prefix, size uint64) (batch *message.Batch) {
	q.requestChan[prefix] <- size
	return <-q.responseChan[prefix]
}

func (q *queue) Input() <-chan *message.Batch {
	return q.responseChan[INPUT]
}

func (q *queue) Output() <-chan *message.Batch {
	return q.responseChan[OUTPUT]
}

func (q *queue) Accept() <-chan *message.Batch {
	return q.responseChan[ACCEPT]
}

func (q *queue) Reject() <-chan *message.Batch {
	return q.responseChan[REJECT]
}

func (q *queue) query(ctx context.Context, prefix uint64) {

	prefixKey := u64ToBytes(prefix)
	shutdown := make(chan bool)
	defer close(q.responseChan[prefix])

shutdown:
	for {
		select {
		case <-shutdown:
			break shutdown
		case size, ok := <-q.requestChan[prefix]:

			if ok {
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
							shutdown <- true
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
			} else {
				break shutdown
			}
		}
	}
}

func (q *queue) shutdown() (err error) {
	q.cancel()
	return q.db.Close()
}

func (h *head) Get(prefix uint64) []byte {
	h.Lock()
	defer h.Unlock()
	return h.m[prefix]
}

func (h *head) Set(prefix uint64, key []byte) {
	h.Lock()
	defer h.Unlock()
	copy(h.m[prefix], key)
}

func (h *head) Reset(prefix uint64) {
	h.Lock()
	defer h.Unlock()
	h.m[prefix] = nil
	h.m[prefix] = make([]byte, 32)
}
