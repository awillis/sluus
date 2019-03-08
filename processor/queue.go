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

var prefixMap map[plugin.Type][]uint64

type (
	queue struct {
		pType        plugin.Type
		numInFlight  uint64
		batchSize    uint64
		batchTimeout time.Duration
		pollInterval time.Duration
		opts         badger.Options
		db           *badger.DB
		head         head
		requestChan  map[uint64]chan bool
		responseChan map[uint64]chan *message.Batch
		logger       *zap.SugaredLogger
	}

	head struct {
		sync.RWMutex
		m map[uint64][]byte
	}
)

func init() {
	prefixMap[plugin.SOURCE] = []uint64{OUTPUT, REJECT, ACCEPT}
	prefixMap[plugin.CONDUIT] = []uint64{OUTPUT, REJECT, ACCEPT, INPUT}
	prefixMap[plugin.SINK] = []uint64{INPUT}
}

func newQueue(pType plugin.Type) (q *queue) {

	dbopts := badger.DefaultOptions
	dbopts.SyncWrites = false

	return &queue{
		pType:        pType,
		requestChan:  make(map[uint64]chan bool),
		responseChan: make(map[uint64]chan *message.Batch),
		opts:         dbopts,
		head: head{
			m: make(map[uint64][]byte),
		},
	}
}

func (q *queue) Initialize() (err error) {

	for _, prefix := range prefixMap[q.pType] {
		q.requestChan[prefix] = make(chan bool, q.numInFlight)
		q.responseChan[prefix] = make(chan *message.Batch, q.numInFlight)
	}

	if e := os.MkdirAll(q.opts.Dir, 0755); e != nil {
		return e
	}
	q.db, err = badger.Open(q.opts)
	return
}

func (q *queue) Start() {
	for _, prefix := range prefixMap[q.pType] {
		go q.query(prefix)
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

func (q *queue) query(prefix uint64) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	ticker := time.NewTicker(q.pollInterval)
	defer ticker.Stop()

	prefixKey := u64ToBytes(prefix)

loop:
	select {
	case <-ticker.C:
		goto loop
	case _, ok := <-q.requestChan[prefix]:

		if ok {
			batch := message.NewBatch(q.batchSize)
			var shutdown bool

			err := q.db.View(func(txn *badger.Txn) (e error) {

				opts := badger.IteratorOptions{
					PrefetchValues: true,
					PrefetchSize:   int(q.batchSize),
				}

				if opts.PrefetchSize > 128 {
					opts.PrefetchSize = 128
				}

				iter := txn.NewIterator(opts)
				defer iter.Close()

				timeout, cancel := context.WithTimeout(context.Background(), q.batchTimeout)
				defer cancel()

			fetch:
				for iter.Seek(q.head.Get(prefix)); iter.ValidForPrefix(prefixKey); iter.Next() {

					select {
					case <-timeout.Done():
						break fetch
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

						if err := batch.AddE(msg); err == message.ErrBatchFull {
							break fetch
						}
					}
				}

				// if there are more records to be read, copy the key to seed the next read
				// otherwise clear the head so that the next read can start at the beginning
				if iter.ValidForPrefix(prefixKey) {
					q.head.Set(prefix, iter.Item().Key())
				} else {
					q.head.Reset(prefix)
				}
				return
			})

			if err != nil {
				q.Logger().Error(errors.WithStack(err))
			}

			if q.pType == plugin.SINK {
				q.Logger().Info("I am a true snake")
			}
			if batch.Count() > 0 && shutdown == false {
				if q.pType == plugin.SINK {
					q.Logger().Infof("queue sink batch: %d", batch.Count())
				}
				q.responseChan[prefix] <- batch
			}
			goto loop
		} else {
			break
		}
	}
}

func (q *queue) shutdown() (err error) {
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
