package tcp

import (
	"bufio"
	"context"
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
	start     chan *net.TCPConn
	end       chan *net.TCPConn
	wg        *sync.WaitGroup
	sock      *net.TCPListener
	conntable *sync.Map
}

func (s *Source) Options() interface{} {
	return s.opts
}

func (s *Source) Initialize() (err error) {
	// validate configuration and set reasonable defaults
	plugin.Validate(s.opts,
		s.opts.defaultPort(),
		s.opts.defaultBatchSize(),
		s.opts.defaultBufferSize(),
		s.opts.defaultSockBufferSize(),
	)

	s.opts.logCurrentConfig(s.Logger())

	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", s.opts.Port))
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

	return
}

func (s *Source) Start(ctx context.Context) {
	go s.Listener()
	go s.Handler(ctx)
	go s.Collector(ctx)
	go s.Closer(ctx)
}

func (s *Source) Produce() <-chan *message.Batch {
	return s.batch
}

func (s *Source) Shutdown() (err error) {
	s.Logger().Error(s.sock.Close())
	s.Logger().Info("socket is closed")
	close(s.end)
	close(s.message)
	s.wg.Wait()
	return
}

func (s *Source) Listener() {

	s.wg.Add(1)
	defer close(s.start)
	defer s.wg.Done()
	defer s.Logger().Info("exit listener")

	for {
		conn, err := s.sock.AcceptTCP()

		if err != nil {
			if strings.Contains(err.Error(), "accept tcp") {
				if err = s.sock.SetDeadline(time.Now().Add(3 * time.Second)); err != nil {
					s.Logger().Info("shutting down tcp listener")
					break
				} else {
					continue
				}
			}
			s.Logger().Errorf("socket error: %v", err)
		}

		if err := conn.SetReadBuffer(int(s.opts.SockBufferSize)); err != nil {
			s.Logger().Errorf("error setting socket buffer size: %v", err)
		}
		s.conntable.Store(conn.RemoteAddr().String(), *conn)
		s.start <- conn
	}
}

func (s *Source) Handler(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()
	ticker := time.NewTicker(time.Duration(s.opts.PollInterval) * time.Millisecond)
	defer ticker.Stop()

loop:
	select {
	case <-ctx.Done():
		break
	case conn := <-s.start:
		go s.handleConnection(conn)
		goto loop
	case <-ticker.C:
		runtime.Gosched()
		goto loop
	}
}

func (s *Source) handleConnection(conn *net.TCPConn) {
	s.wg.Add(1)
	defer s.wg.Done()
	defer s.Logger().Info("exit connection handler")
	s.Logger().Info("start connection handler")

	scanner := bufio.NewScanner(bufio.NewReader(conn))
	scanner.Buffer(make([]byte, 0, s.opts.BufferSize), int(s.opts.BufferSize))

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

func (s *Source) Collector(ctx context.Context) {
	s.wg.Add(1)
	defer close(s.batch)
	defer s.wg.Done()
	defer s.Logger().Info("exit collector")
	s.Logger().Info("start collector")

	batch := message.NewBatch(s.opts.BatchSize)
	ticker := time.NewTicker(time.Duration(s.opts.PollInterval) * time.Millisecond)
	defer ticker.Stop()

loop:
	select {
	case <-ctx.Done():
		break
	case msg := <-s.message:
		if batch.Count() < s.opts.BatchSize {
			if err := batch.Add(msg); err != nil {
				s.Logger().Error(err)
			}
		} else {
			b := message.NewBatch(batch.Count())
			for msg := range batch.Iter() {
				_ = b.Add(msg)
			}
			s.batch <- b
			batch.Clear()
		}
		goto loop
	case <-ticker.C:
		runtime.Gosched()
		goto loop
	}
}

func (s *Source) Closer(ctx context.Context) {
	s.wg.Add(1)
	defer s.wg.Done()
	ticker := time.NewTicker(time.Duration(s.opts.PollInterval) * time.Millisecond)
	defer ticker.Stop()

loop:
	select {
	case <-ctx.Done():
		break
	case conn := <-s.end:
		if err := conn.Close(); err != nil {
			s.Logger().Errorf("error while closing connection: %v", err)
		} else {
			s.conntable.Delete(conn.RemoteAddr().String())
		}
		goto loop
	case <-ticker.C:
		runtime.Gosched()
		goto loop
	}
}
