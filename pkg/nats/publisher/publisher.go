package publisher

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
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
	natsConfig := config.Get().Nats

	// Connect to a server
	logger.Log.Debugf("Connecting to nats@%s", natsConfig.URL)

	nc, err := nats.Connect(natsConfig.URL)
	if err != nil {
		return nil, err
	}

	Gibson := &Gibson{
		conn:  nc,
		input: make(chan []byte),
		stats: newStats(),
		topic: natsConfig.Publisher.Topic,
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
