package processor

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/binary"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/dgraph-io/badger"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
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
		batchTimeout time.Duration
		pollInterval time.Duration
		opts         badger.Options
		db           *badger.DB
		wg           *sync.WaitGroup
		head         head
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

	dbopts := badger.DefaultOptions
	dbopts.SyncWrites = false

	rqChan := make(map[uint64]chan uint64)
	rsChan := make(map[uint64]chan *message.Batch)

	for _, prefix := range []uint64{INPUT, OUTPUT, REJECT, ACCEPT} {
		rqChan[prefix] = make(chan uint64)
		rsChan[prefix] = make(chan *message.Batch)
	}

	return &queue{
		pType:        pType,
		wg:           new(sync.WaitGroup),
		requestChan:  rqChan,
		responseChan: rsChan,
		opts:         dbopts,
		head: head{
			m: make(map[uint64][]byte),
		},
	}
}

func (q *queue) Initialize() (err error) {

	if e := os.MkdirAll(q.opts.Dir, 0755); e != nil {
		return e
	}
	q.db, err = badger.Open(q.opts)
	return
}

func (q *queue) Start(ctx context.Context) {

	switch q.pType {
	case plugin.SOURCE:
		go q.query(ctx, OUTPUT)
		go q.query(ctx, ACCEPT)
		go q.query(ctx, REJECT)
	case plugin.CONDUIT:
		go q.query(ctx, INPUT)
		go q.query(ctx, OUTPUT)
		go q.query(ctx, ACCEPT)
		go q.query(ctx, REJECT)
	case plugin.SINK:
		go q.query(ctx, INPUT)
	}
}

func (q *queue) Logger() *zap.SugaredLogger {
	return q.logger.With("component", "queue")
}

func u64ToBytes(i uint64) (b []byte) {
	b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return
}

func (q *queue) Put(prefix uint64, batch *message.Batch) {

	err := q.db.Update(func(txn *badger.Txn) (e error) {

		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()

		hash := md5.New()
		prefixKey := u64ToBytes(prefix)

		key := make([]byte, 0, 24+md5.Size)
		key = append(key, prefixKey...)

		for msg := range batch.Iter() {

			payload, err := msg.ToBytes()

			if err != nil {
				q.Logger().Error(err)
			}

			timeKey := u64ToBytes(uint64(time.Now().UnixNano()))
			sizeKey := u64ToBytes(uint64(len(payload)))
			contentKey := hash.Sum(payload)

			key = append(key, sizeKey...)
			key = append(key, timeKey...)
			key = append(key, contentKey...)

			e = txn.Set(key, payload)
			key = key[:8]
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
		q.Logger().Error(errors.WithStack(err))
	}

	return
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

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	q.wg.Add(1)
	defer q.wg.Done()

	ticker := time.NewTicker(q.pollInterval)
	defer ticker.Stop()

	prefixKey := u64ToBytes(prefix)

loop:
	select {
	case <-ctx.Done():
		break
	case <-ticker.C:
		runtime.Gosched()
		goto loop
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

				timeout, cancel := context.WithTimeout(ctx, q.batchTimeout)
				defer cancel()

			timeout:
				for iter.Seek(q.head.Get(prefix)); iter.ValidForPrefix(prefixKey); iter.Next() {

					select {
					case <-timeout.Done():
						break timeout
					default:
						item := iter.Item()
						value, err := item.Value()

						if err != nil {
							e = err
						}

						msg, err := message.FromBytes(value)

						if err != nil {
							e = err
						}

						if err := batch.Add(msg); err != nil {
							break
						}
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

			if batch.Count() > 0 {
				q.responseChan[prefix] <- batch
			}

			if err != nil {
				q.Logger().Error(errors.WithStack(err))
			}
		}
		runtime.Gosched()
		goto loop
	}

}

func (q *queue) shutdown() (err error) {
	q.wg.Wait()
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
	h.m[prefix] = make([]byte, 0, 32)
}
