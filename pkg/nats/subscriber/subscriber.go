package subscriber

import (
	"context"
	"reflect"

	"github.com/nats-io/nats.go"
	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
	natsGibson "github.com/simplefxn/go-gibson/pkg/nats"
)

// Gibson structure
type Gibson struct {
	conn   *nats.Conn
	stats  *Stats
	topics []Topic
}

type Topic struct {
	callback func(msg *nats.Msg)
	name     string
	chann    chan *nats.Msg
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
		conn:  nc,
		stats: newStats(),
	}

	return Gibson, nil
}

func (g *Gibson) Add(topic string, callback func(msg *nats.Msg)) {
	t := Topic{
		name:     topic,
		callback: callback,
		chann:    make(chan *nats.Msg),
	}
	g.topics = append(g.topics, t)
}

// Run main loop for the receiver , call the callback for every message
func (g *Gibson) Run(ctx context.Context) error {

	go func() {
		cases := make([]reflect.SelectCase, len(g.topics))
		for i, t := range g.topics {
			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(t.chann)}
		}

		for {
			i, value, ok := reflect.Select(cases)
			if ok {
				msg := value.Interface().(nats.Msg)
				g.topics[i].callback(&msg)
			}

		}

	}()

	go func() {
		<-ctx.Done()
		g.conn.Close()

	}()

	for i := range g.topics {
		// Simple Async Subscriber
		g.conn.Subscribe(g.topics[i].name, func(m *nats.Msg) {
			// Apply stats logic
			g.topics[i].chann <- m
		})
	}

	return nil
}
