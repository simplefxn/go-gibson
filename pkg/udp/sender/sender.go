package sender

import (
	"context"
	"fmt"
	"net"

	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
)

// Gibson structure
type Gibson struct {
	conn  *net.UDPConn
	input chan []byte
	stats *Stats
}

// New creates a new UDP sender
func New() (*Gibson, error) {
	conf := config.Get()

	addr := fmt.Sprintf("%s:%d", conf.Receiver.Address, conf.Receiver.Port)

	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	c, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return nil, err
	}

	logger.Log.Infof("Remote UDP address : %s \n", c.RemoteAddr().String())
	logger.Log.Infof("Local UDP client address : %s \n", c.LocalAddr().String())

	Gibson := &Gibson{
		conn:  c,
		input: make(chan []byte),
		stats: newStats(),
	}

	return Gibson, nil
}

// Close closes the udp channel
func (u *Gibson) Close() error {
	return u.conn.Close()
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
			case pkt := <-u.input:
				size, err := u.conn.Write(pkt)
				if err != nil {
					u.stats.IncErr()
				}
				u.stats.IncMsg(size)
			}
		}
	}()
	return nil
}
