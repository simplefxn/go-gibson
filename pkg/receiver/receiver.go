package receiver

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
	"github.com/spf13/cobra"
)

// Receiver struct
type Gibson struct {
	conn     *net.UDPConn
	stats    *Stats
	callback func(msg []byte) error
}

func New(cmd *cobra.Command, callback func(msg []byte) error) (*Gibson, error) {
	conf := config.Get()

	sAddr := fmt.Sprintf("%s:%d", conf.Receiver.Address, conf.Receiver.Port)
	s, err := net.ResolveUDPAddr("udp4", sAddr)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp4", s)
	if err != nil {
		return nil, err
	}
	stats := newStats()

	u := &Gibson{
		conn:     conn,
		stats:    stats,
		callback: callback,
	}

	return u, nil
}

func (g *Gibson) read(chan []byte) error {
	buffer := make([]byte, 2048)
	for {
		n, addr, err := g.conn.ReadFromUDP(buffer)
		logger.Log.Debugf("received %d bytes from %s", n, addr.String())
		if err != nil {
			g.stats.IncErr()
			return err
		}
		g.stats.IncMsg(n)
	}
}

// Run main loop for the receiver , call the callback for every message
func (g *Gibson) Run(ctx context.Context) error {
	msgChannel := make(chan []byte)

	// Read
	go func() {
		retries := 5
		timeout := 10 * time.Second
		for {
			err := g.read(msgChannel)
			if err != nil {
				logger.Log.Error(err)
				retries--
			}
			if ctx.Err() != nil || retries <= 0 {
				return
			}
			time.Sleep(timeout)
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				g.conn.Close()
				return
			case pkg := <-msgChannel:
				g.callback(pkg)
			}
		}
	}()
	return nil
}
