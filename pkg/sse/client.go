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
)

//SSE name constants
const (
	eName = "event"
	dName = "data"
)

type Event struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

var (
	//ErrNilChan will be returned by Notify if it is passed a nil channel
	ErrNilChan = fmt.Errorf("nil channel given")
)

type Client struct {
	client *http.Client
	ctx    context.Context
}

func New(ctx context.Context) *Client {
	c := &Client{
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
		ctx: ctx,
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

//SSEvent is a go representation of an http server-sent event
type SSEvent struct {
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

func getEvent(br *bufio.Reader) (*Event, error) {
	delim := []byte{':', ' '}
	currEvent := &Event{}

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

func getEvents(ctx context.Context, br *bufio.Reader, evCh chan<- *Event) error {

	for {
		currEvent, err := getEvent(br)
		if err != nil {
			return err
		}
		if ctx.Err() != nil {
			return nil
		}

		// Increment internal metrics counter
		eventCounter.Inc()
		evCh <- currEvent

	}
}

func (c *Client) Start(url string, callback func(*Event)) {
	var wg sync.WaitGroup
	// Make a receive channel for getting messages from the http response
	recvChan := make(chan *Event)
	ctxDone := false

	// Connection goroutine, connect, fecth events , repeat
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if ctxDone {
				logger.Log.Info("Context close, exiting sse")
				return
			}
			logger.Log.Info("Connecting to SSE Server")
			res, err := c.clientConnect(url)
			if err != nil {
				logger.Log.Info("Client connect skip until next cycle.")
				continue
			}

			// GoRoutine that will listen for the context and close the response if the context
			// is closed
			go func(ctx context.Context, res *http.Response) {
				<-ctx.Done()
				ctxDone = true
				logger.Log.Info("Closing service side response")
				res.Body.Close()
			}(c.ctx, res)

			// Create bufio reader
			br := bufio.NewReader(res.Body)
			// Loop for all events and send them to the recv Channel
			// this blocks until the response is close
			err = getEvents(c.ctx, br, recvChan)
			// If the goRoutine context is dome
			if err != nil {
				SSERestartCounter.Inc()
				logger.Log.Info("Error from getEvents due to %s", err.Error())
				//res.Body.Close()
				continue
			}
		}
	}()

	// Report goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		ticker := time.NewTicker(config.Get().Globals.ReportInterval)
	report:
		for {
			select {
			case <-c.ctx.Done():
				logger.Log.Info("Reporting received context closing, stoping now")
				//close(recvChan) TODO: This channel is leaking
				ticker.Stop()
				break report
			case <-ticker.C:
				x := eventCounter.Load()
				logger.Log.Infof("Rate %s msg/int", strconv.FormatUint(x, 10))
				eventCounter.Store(0)
			}
		}
	}()

	// SSE reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
	reader:
		for {
			select {
			case <-c.ctx.Done():
				logger.Log.Info("SSE reader received context closing, stoping now")
				//close(recvChan) TODO: This channel is leaking
				break reader
			// If we receive a message, call back to user function
			case msg := <-recvChan:
				callback(msg)
			}

		}
	}()

	wg.Wait()
}
