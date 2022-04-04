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
	stats  *Stats
	topics map[string]Topic
}

type Topic struct {
	conn  *nats.Conn
	cb    func(msg *nats.Msg)
	chann chan *nats.Msg
}

// New creates a new UDP sender
func New() (*Gibson, error) {
	Gibson := &Gibson{
		stats:  newStats(),
		topics: make(map[string]Topic),
	}

	return Gibson, nil
}

func (g *Gibson) Add(topic string, callback func(msg *nats.Msg)) error {
	natsConfig := config.Get().Nats
	natsTLSconfig := natsGibson.CreateTlsConfiguration(natsConfig)
	// Connect to a server
	logger.Log.Debugf("Connecting to nats@%s for %s", natsConfig.URL, topic)

	nc, err := nats.Connect(natsConfig.URL, nats.Secure(natsTLSconfig))
	if err != nil {
		return err
	}

	t := Topic{
		conn:  nc,
		cb:    callback,
		chann: make(chan *nats.Msg),
	}
	g.topics[topic] = t

	return nil
}

// Run main loop for the receiver , call the callback for every message
func (g *Gibson) Run(ctx context.Context) error {
	/* TODO: This code is cleaner , but it does not work, it always return the index of the last topic
	            entered in the list
	   	go func() {
	   		cases := make([]reflect.SelectCase, len(g.topics))
	   		for i, t := range g.topics {
	   			logger.Log.Debugf("Listening in %s", t.name)
	   			cases[i] = reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(t.chann)}
	   		}

	   		for {
	   			i, value, ok := reflect.Select(cases)
	   			logger.Log.Debugf("Received msg on channel %d", i)
	   			if ok {
	   				msg := value.Interface().(*nats.Msg)
	   				g.topics[i].cb(msg)
	   			}

	   		}
	   	}()
	*/

	for subject, topic := range g.topics {
		go func(subject string, topic Topic) {
			logger.Log.Debugf("Subscribing to subject %s", subject)
			topic.conn.Subscribe(subject, func(m *nats.Msg) {
				topic.chann <- m
			})
			for {
				select {
				case <-ctx.Done():
					close(topic.chann)
					topic.conn.Close()
					return
				case msg := <-topic.chann:
					topic.cb(msg)
				}
			}
		}(subject, topic)
	}
	return nil
}
