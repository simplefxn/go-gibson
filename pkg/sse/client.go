package sse // import "astuart.co/go-sse"

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
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

// Event structure for SSE
type Event struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

var (
	//ErrNilChan will be returned by Notify if it is passed a nil channel
	ErrNilChan = fmt.Errorf("nil channel given")
)

// Client for SSE
type Client struct {
	client *http.Client
}

// New creates a new SSE client
func New() *Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.IdleConnTimeout = config.Get().SSE.Section.IdleTimeout

	c := &Client{
		client: &http.Client{
			// Timeout:   5 * time.Minute, TODO: We dont need to close the conection , just when is idle via the transport
			Transport: t,
		},
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
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recovering from panic in getEvents error is: %v \n", r)
		}
	}()

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

func (c *Client) Start(ctx context.Context, url string, callback func(*Event)) {
	var wg sync.WaitGroup
	// Make a receive channel for getting messages from the http response
	recvChan := make(chan *Event)

	// Connection goroutine, connect, fecth events , repeat
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			logger.Log.Info("Exiting goroutine 0")
		}()
		for {
			if ctx.Err() != nil {
				logger.Log.Info("Context close, exiting sse")
				return
			}
			logger.Log.Info("Connecting to SSE Server")
			res, err := c.clientConnect(url)
			if err != nil {
				logger.Log.Info("Client connect skip until next cycle.")
				continue
			}

			// Create bufio reader
			br := bufio.NewReader(res.Body)
			// Loop for all events and send them to the recv Channel
			// this blocks until the response is close
			err = getEvents(ctx, br, recvChan)
			// If the goRoutine context is done
			if err != nil {
				SSERestartCounter.Inc()
				logger.Log.Info("Error from getEvents due to %s", err.Error())
				res.Body.Close()
				continue
			}
		}
	}()

	// Report goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			logger.Log.Info("Exiting goroutine 1")
		}()
		ticker := time.NewTicker(config.Get().SSE.Section.ReportInterval)
		minuteTicker := time.NewTicker(time.Minute * 1)
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("Report routine received context closing, stoping now")
				close(recvChan) // TODO: This channel is leaking
				ticker.Stop()
				return
			case <-ticker.C:
				x := eventCounter.Load()
				measurementsMutex.Lock()
				measurements = append(measurements, x)
				measurementsMutex.Unlock()
				// reset counter
				eventCounter.Store(0)
			case <-minuteTicker.C:
				measurementsMutex.Lock()
				var total uint64
				for _, n := range measurements {
					total = total + n
				}
				avg := total / uint64(len(measurements))
				logger.Log.Infof("SSE Avg Rate %d msg/min", avg)
				measurements = nil
				measurementsMutex.Unlock()
			}
		}
	}()

	// SSE reader goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		defer func() {
			logger.Log.Info("Exiting goroutine 2")
		}()
	reader:
		for {
			select {
			case <-ctx.Done():
				logger.Log.Info("SSE reader received context closing, stoping now")
				close(recvChan) //TODO: This channel is leaking
				break reader
			// If we receive a message, call back to user function
			case msg := <-recvChan:
				callback(msg)
			}

		}
	}()

	wg.Wait()
	logger.Log.Infof("Out of SSE client")
}
