package sse

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/atomic"
)

var eventCounter atomic.Uint64
var measurements []uint64
var measurementsMutex sync.RWMutex

var SSERestartCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "gibson_sse_restart_count",
	Help: "The number of times the client SSE has restarted",
})
