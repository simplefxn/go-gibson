package sse // import "astuart.co/go-sse"

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
	"github.com/simplefxn/go-gibson/pkg/metrics"
)

//SSE name constants
const (
	eName = "event"
	dName = "data"
)

type SSE_RIS struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

var (
	//ErrNilChan will be returned by Notify if it is passed a nil channel
	ErrNilChan = fmt.Errorf("nil channel given")
)

type Client struct {
	client *http.Client
	uri    string
	conf   config.Service
	ctx    context.Context
}

func New(ctx context.Context, conf config.Service, uri string) *Client {
	c := &Client{
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
		conf: conf,
		uri:  uri,
		ctx:  ctx,
	}

	return c
}

func liveReq(verb, uri string, body io.Reader) (*http.Request, error) {
	req, err := GetReq(verb, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "text/event-stream")

	return req, nil
}

//Event is a go representation of an http server-sent event
type Event struct {
	Type string
	Data io.Reader
}

//GetReq is a function to return a single request. It will be used by notify to
//get a request and can be replaces if additional configuration is desired on
//the request. The "Accept" header will necessarily be overwritten.
var GetReq = func(verb, uri string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(verb, uri, body)
}

func (c *Client) clientConnect(uri string) (*http.Response, error) {

	req, err := liveReq("GET", uri, nil)
	if err != nil {
		return nil, fmt.Errorf("error getting sse request: %v", err)
	}

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error performing request for %s: %v", uri, err)
	}
	return res, nil
}

func getEvent(br *bufio.Reader) (*SSE_RIS, error) {
	delim := []byte{':', ' '}
	currEvent := &SSE_RIS{}

	for {
		bs, err := br.ReadBytes('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if len(bs) < 2 {
			continue
		}

		spl := bytes.Split(bs, delim)

		if len(spl) < 2 {
			continue
		}

		switch string(spl[0]) {
		case eName:
			currEvent.Type = string(bytes.TrimSpace(spl[1]))
		case dName:
			currEvent.Data = string(bytes.TrimSpace(spl[1]))
			return currEvent, nil
		}
		if err == io.EOF {
			return nil, err
		}
	}
}

func getEvents(br *bufio.Reader, evCh chan<- *SSE_RIS) error {

	for {
		currEvent, err := getEvent(br)
		if err != nil {
			logger.Log.Errorf("Error getting event: %s", err.Error())
			return err
		}
		// Increment internal metrics counter
		metrics.SSEEventCounter.Inc()
		evCh <- currEvent
	}
}

func (c *Client) Start(cb func(*SSE_RIS)) {
	var wg sync.WaitGroup
	// Make a receive channel for getting messages from the http response
	recvChan := make(chan *SSE_RIS)
	ctxDone := false

	// Main goroutine, connect, fecth event , repeat
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if ctxDone {
				return
			}
			res, err := c.clientConnect(c.uri)
			if err != nil {
				logger.Log.Info("Client connect skip until next cycle.")
				continue
			}

			// GoRoutine that will listen for the context and close the response if the context
			// is closed
			go func(ctx context.Context, res *http.Response) {
				<-ctx.Done()
				ctxDone = true
				logger.Log.Info("Received context close, closing service side response")
				res.Body.Close()
			}(c.ctx, res)

			// Create bufio reader
			br := bufio.NewReader(res.Body)
			// Loop for all events and send them to the recv Channel
			// this blocks until the response is close
			err = getEvents(br, recvChan)
			// If the goRoutine context is dome
			if err != nil {
				logger.Log.Info("Error from getting events from connection, skip until next cycle")
				res.Body.Close()
				continue
			}
		}
	}()

	// Goroutine to display in log the rate
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Log.Infof("Starting producer monitoring every %f", c.conf.MonInverval.Seconds())

	ticker_outside:
		for {
			select {

			case <-c.ctx.Done():
				logger.Log.Info("SSE client receive signal to stop, closing receive channel")
				close(recvChan)

				break ticker_outside
			}
		}
		logger.Log.Info("Exit producer monitiring goroutine")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(c.conf.MonInverval)
	outside_cb:
		for {
			select {
			case <-c.ctx.Done():
				logger.Log.Info("SSE client receive signal to stop, closing receive channel")
				close(recvChan)
				ticker.Stop()
				break outside_cb
			// If we receive a message, call back to user function
			case msg := <-recvChan:
				cb(msg)
			case <-ticker.C:
				x := metrics.SSEEventCounter.Load()
				logger.Log.Infof("Rate %s msg/int", strconv.FormatUint(x, 10))
				metrics.SSEEventCounter.Store(0)
			}

		}
	}()

	wg.Wait()
}
