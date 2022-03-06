package publisher

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/simplefxn/go-gibson/pkg/config"
)

// Gibson structure
type Gibson struct {
	topic string
	conn  *nats.Conn
	input chan []byte
	stats *Stats
}

// New creates a new UDP sender
func New() (*Gibson, error) {
	conf := config.Get()

	// Connect to a server
	nc, _ := nats.Connect(conf.Nats.URL)

	Gibson := &Gibson{
		conn:  nc,
		input: make(chan []byte),
		stats: newStats(),
		topic: config.Nats.Publisher.Topic,
	}

	return Gibson, nil
}

// Close closes the udp channel
func (u *Gibson) Close() {
	u.conn.Close()
}

// Input returns a channel to send messages
func (u *Gibson) Input() chan<- []byte {
	return u.input
}

// Run execute the non-blocking main loop
func (u *Gibson) Run(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-u.input:
				u.conn.Publish(u.topic, msg)
				u.stats.IncMsg(len(msg))
			}
		}
	}()
	return nil
}
