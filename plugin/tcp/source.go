package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/awillis/sluus/message"
	"github.com/awillis/sluus/plugin"
)

var _ plugin.Interface = new(Source)
var _ plugin.Producer = new(Source)

var (
	ErrSocketConfig = errors.New("unable to configure socket")
)

type Source struct {
	plugin.Base
	opts      *options
	batch     chan *message.Batch
	message   chan *message.Message
	wg        *sync.WaitGroup
	sock      *net.TCPListener
	start     chan *net.TCPConn
	end       chan *net.TCPConn
	conntable *sync.Map
}

func (s *Source) Options() interface{} {
	return s.opts
}

func (s *Source) Initialize() (err error) {
	// validate configuration and set reasonable defaults
	plugin.Validate(s.opts, s.opts.Port(), s.opts.BatchSize(), s.opts.BufferSize(), s.opts.SockBufferSize())

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", s.opts.port))
	if err != nil {
		return errors.Wrap(ErrSocketConfig, err.Error())
	}

	s.sock, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return errors.Wrap(ErrSocketConfig, err.Error())
	}

	if err = s.sock.SetDeadline(time.Now().Add(10 * time.Second)); err != nil {
		return errors.Wrap(ErrSocketConfig, err.Error())
	}

	go s.Listener()
	for i := 0; i < runtime.NumCPU(); i++ {
		go s.Handler()
		go s.Collector()
	}
	go s.Closer()

	s.Logger().Info("initialized")
	return
}

func (s *Source) Produce() (ch chan *message.Batch) {
	return s.batch
}

func (s *Source) Shutdown() (err error) {

	if s.sock != nil {
		err = s.sock.Close()
	}

	close(s.message)
	close(s.batch)
	s.wg.Wait()
	return
}

func (s *Source) Listener() {

	s.wg.Add(1)
	defer close(s.start)
	defer s.wg.Done()

	for {
		conn, err := s.sock.AcceptTCP()

		if err != nil {
			if strings.Contains(err.Error(), "accept tcp") {
				if err = s.sock.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
					s.Logger().Info("shutting down tcp acceptor")
					break
				} else {
					continue
				}
			}
			s.Logger().Errorf("socket error: %v", err)
		}

		if err := conn.SetReadBuffer(s.opts.sockBufferSize); err != nil {
			s.Logger().Errorf("error setting socket buffer size: %v", err)
		}
		s.conntable.Store(conn.RemoteAddr().String(), *conn)
		s.start <- conn
	}
}

func (s *Source) Handler() {

	s.wg.Add(1)
	defer close(s.end)
	defer s.wg.Done()

shutdown:
	for {
		select {
		case conn, ok := <-s.start:
			if ok {
				go s.handleConnection(conn)
			} else {
				break shutdown
			}
		}
	}
}

func (s *Source) handleConnection(conn *net.TCPConn) {
	s.wg.Add(1)
	defer s.wg.Done()

	scanner := bufio.NewScanner(bufio.NewReader(conn))
	scanner.Buffer(make([]byte, 0, s.opts.bufferSize), s.opts.bufferSize)

	for {
		if scanner.Scan() {

			msg, err := message.WithContent(json.RawMessage(scanner.Text()))

			if err != nil {
				s.Logger().Error(err)
			}

			msg.Direction = message.Message_PASS
			msg.MarkReceived()
			s.message <- msg
		} else {
			if err := scanner.Err(); err != nil {
				s.Logger().Errorf("error while reading from client: %v", err)
			}
			break
		}
	}
	s.end <- conn
}

func (s *Source) Collector() {

	batch := message.NewBatch(s.opts.batchSize)

	for {
		select {
		case msg, ok := <-s.message:
			if !ok {
				break
			} else {
				if batch.Count() < uint64(s.opts.batchSize) {
					if err := batch.Add(msg); err != nil {
						s.Logger().Error(err)
					}
				} else {
					b := message.NewBatch(int(batch.Count()))
					for msg := range batch.Iter() {
						_ = b.Add(msg)
					}
					s.batch <- b
					batch.Clear()
				}
			}
		}
	}
}

func (s *Source) Closer() {

	s.wg.Add(1)
	defer s.wg.Done()

shutdown:
	for {
		select {
		case conn, ok := <-s.end:
			if ok {
				if err := conn.Close(); err != nil {
					s.Logger().Errorf("error while closing connection: %v", err)
				} else {
					s.conntable.Delete(conn.RemoteAddr().String())
				}
			} else {
				break shutdown
			}
		}
	}
}
