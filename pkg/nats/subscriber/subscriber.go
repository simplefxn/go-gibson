package subscriber

import (
	"context"

	"github.com/nats-io/nats.go"
	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
	natsGibson "github.com/simplefxn/go-gibson/pkg/nats"
)

// Gibson structure
type Gibson struct {
	topic    string
	conn     *nats.Conn
	input    chan []byte
	stats    *Stats
	callback func(msg *nats.Msg)
}

// New creates a new UDP sender
func New(topic string, callback func(msg *nats.Msg)) (*Gibson, error) {
	natsConfig := config.Get().Nats

	natsTLSconfig := natsGibson.CreateTlsConfiguration(natsConfig)

	// Connect to a server
	logger.Log.Debugf("Connecting to nats@%s", natsConfig.URL)

	nc, err := nats.Connect(natsConfig.URL, nats.Secure(natsTLSconfig))
	if err != nil {
		return nil, err
	}

	Gibson := &Gibson{
		conn:     nc,
		input:    make(chan []byte),
		stats:    newStats(),
		topic:    topic,
		callback: callback,
	}

	return Gibson, nil
}

// Run main loop for the receiver , call the callback for every message
func (g *Gibson) Run(ctx context.Context) error {

	msgChannel := make(chan *nats.Msg)

	go func() {
		for {
			select {
			case <-ctx.Done():
				g.conn.Close()
				return
			case msg := <-msgChannel:
				g.callback(msg)
			}
		}
	}()

	// Simple Async Subscriber
	g.conn.Subscribe(g.topic, func(m *nats.Msg) {
		// Apply stats logic
		msgChannel <- m
	})

	return nil
}
