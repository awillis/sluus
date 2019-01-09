package kafka

import (
	"sync"

	"github.com/segmentio/kafka-go"

	"github.com/awillis/sluus/plugin"
)

type Sink struct {
	plugin.Base
	msgs   chan []byte
	wg     *sync.WaitGroup
	writer *kafka.Writer
	opt    *options
}

type options struct {

	// Bootstrap Server ( "host:port" )
	BootstrapServer string `mapstructure:"bootstrap_server"`
	// Broker list
	Brokers []string `mapstructure:"brokers"`
	// Kafka topic
	TopicID string `mapstructure:"topic_id" validate:"required"`
	// Kafka client id
	ClientID string `mapstructure:"client_id"`
	// Balancer ( roundrobin, hash or leastbytes )
	Balancer string `mapstructure:"balancer"`
	// Compression algorithm ( 'gzip', 'snappy', or 'lz4' )
	Compression string `mapstructure:"compression"`
	// Max Attempts
	MaxAttempts int `mapstructure:"max_attempts"`
	// Queue Size
	QueueSize int `mapstructure:"queue_size"`
	// Batch Size
	BatchSize int `mapstructure:"batch_size"`
	// Keep Alive ( in seconds )
	KeepAlive int `mapstructure:"keepalive"`
	// IO Timeout ( in seconds )
	IOTimeout int `mapstructure:"io_timeout"`
	// Required Acks ( number of replicas that must acknowledge write. -1 for all replicas )
	RequiredAcks int `mapstructure:"acks"`
	// Periodic Flush ( length of time in seconds a partially written buffer will live before being flushed )
	PeriodicFlush int `mapstructure:"pflush"`
}

func (s *Sink) Initialize() error {
	var err error
	return err
}

func (s *Sink) Execute() error {
	var err error
	return err
}

func (s *Sink) Shutdown() error {
	var err error
	return err
}
