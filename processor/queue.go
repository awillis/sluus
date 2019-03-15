package processor

import (
	"bytes"
	"context"
	"encoding/binary"
	"hash/crc64"
	"os"
	"runtime"
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

var (
	compass = make(map[plugin.Type][]uint64)
	IN      = []uint64{INPUT}
	OUT     = []uint64{OUTPUT, REJECT, ACCEPT}
)

type (
	queue struct {
		db           *badger.DB
		opts         badger.Options
		head         head
		requestChan  map[uint64]chan bool
		responseChan map[uint64]chan *message.Batch
		wg           *workGroup
		cnf          *config
	}
)

func init() {
	compass[plugin.SOURCE] = OUT
	compass[plugin.CONDUIT] = append(IN, OUT...)
	compass[plugin.SINK] = IN
}

func newQueue(cnf *config) (q *queue) {

	dbopts := badger.DefaultOptions
	dbopts.SyncWrites = false

	return &queue{
		cnf:          cnf,
		requestChan:  make(map[uint64]chan bool),
		responseChan: make(map[uint64]chan *message.Batch),
		opts:         dbopts,
		wg:           new(workGroup),
		head: head{
			m: make(map[uint64][]byte),
		},
	}
}

func (q *queue) Initialize() (err error) {

	for _, direction := range compass[q.cnf.pluginType] {
		q.head.Reset(direction)
		q.requestChan[direction] = make(chan bool, q.cnf.qqRequests)
		q.responseChan[direction] = make(chan *message.Batch, q.cnf.qqRequests)
	}

	if e := os.MkdirAll(q.opts.Dir, 0755); e != nil {
		return e
	}

	q.db, err = badger.Open(q.opts)
	return
}

func (q *queue) Start() {

	ctx, cancel := context.WithCancel(context.Background())
	q.wg.cancel = cancel

	for _, prefix := range compass[q.cnf.pluginType] {
		go q.query(ctx, prefix)
	}
}

func (q *queue) Logger() *zap.SugaredLogger {
	return q.cnf.logger.With("component", "queue")
}

func u64ToBytes(i uint64) (b []byte) {
	b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, i)
	return
}

func (q *queue) Put(direction uint64, batch *message.Batch) {

	err := q.db.Update(func(txn *badger.Txn) (e error) {

		iter := txn.NewIterator(badger.DefaultIteratorOptions)
		defer iter.Close()

		if len(q.head.Get(direction)) > 0 {
			iter.Seek(q.head.Get(direction))
		}

		prefixKey := u64ToBytes(direction)

		key := make([]byte, 0, 32)
		key = append(key, prefixKey...)
		table := crc64.MakeTable(crc64.ECMA)

		for msg := range batch.Iter() {

			payload, err := msg.ToBytes()

			if err != nil {
				q.Logger().Error(err)
			}

			timeKey := u64ToBytes(uint64(time.Now().UnixNano()))
			sizeKey := u64ToBytes(uint64(len(payload)))
			crcKey := u64ToBytes(crc64.Checksum(payload, table))

			key = append(key, sizeKey...)
			key = append(key, timeKey...)
			key = append(key, crcKey...)

			e = txn.Set(key, payload)
			key = key[:8]
		}

		if len(q.head.Get(direction)) > 0 {
			// if the read head is set, remove data
			// from the beginning up to the read head
			for iter.Rewind(); iter.ValidForPrefix(prefixKey); iter.Next() {
				key := iter.Item().Key()
				if bytes.Equal(key, q.head.Get(direction)) {
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

func (q *queue) query(ctx context.Context, direction uint64) {

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	q.wg.Add(1)
	defer q.wg.Done()

	prefixKey := u64ToBytes(direction)

loop:
	select {
	case <-ctx.Done():
		break loop
	default:
		batch := message.NewBatch(q.cnf.batchSize)

		err := q.db.View(func(txn *badger.Txn) (e error) {

			opts := badger.IteratorOptions{
				PrefetchValues: true,
				PrefetchSize:   int(q.cnf.batchSize),
			}

			if opts.PrefetchSize > 128 {
				opts.PrefetchSize = 128
			}

			iter := txn.NewIterator(opts)
			defer iter.Close()

			timeout, cancel := context.WithTimeout(ctx, q.cnf.batchTimeout)
			defer cancel()

		fetch:
			for iter.Seek(q.head.Get(direction)); iter.ValidForPrefix(prefixKey); iter.Next() {

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
				key := iter.Item().KeyCopy(nil)
				q.head.Set(direction, key)
			} else {
				q.head.Reset(direction)
			}
			return
		})

		if err != nil {
			q.Logger().Error(errors.WithStack(err))
		}

		if batch.Count() > 0 {
			q.responseChan[direction] <- batch
		} else {
			time.Sleep(q.cnf.pollInterval)
		}
		goto loop
	}
}

func (q *queue) shutdown() (err error) {
	q.Logger().Info("queue query shutdown")
	q.wg.Shutdown()
	q.Logger().Info("queue db close")
	return q.db.Close()
}
