package receiver

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/atomic"
)

type Stats struct {
	totalMsg       atomic.Uint64
	totalBytes     atomic.Uint64
	totalErr       atomic.Uint64
	promTotalBytes prometheus.Counter
	promTotalMsgs  prometheus.Counter
	promErrMsg     prometheus.Counter
}

func newStats() *Stats {
	// Initialize metrics counter
	s := &Stats{
		promTotalBytes: promauto.NewCounter(prometheus.CounterOpts{
			Name: "ugibson_total_bytes_received",
			Help: "The total number of bytes received",
		}),
		promTotalMsgs: promauto.NewCounter(prometheus.CounterOpts{
			Name: "ugibson_total_msgs_received",
			Help: "The total number of msgs received",
		}),
		promErrMsg: promauto.NewCounter(prometheus.CounterOpts{
			Name: "ugibson_total_msg_errors",
			Help: "The total number of msgs errors",
		}),
	}
	s.totalMsg.Store(0)
	s.totalBytes.Store(0)
	s.totalErr.Store(0)
	return s
}

func (s *Stats) GetTotalMsg() uint64 {
	return s.totalMsg.Load()
}

func (s *Stats) GetTotalBytes() uint64 {
	return s.totalBytes.Load()
}

func (s *Stats) GetTotalErr() uint64 {
	return s.totalErr.Load()
}

func (s *Stats) IncMsg(size int) {
	s.totalMsg.Inc()
	s.totalBytes.Add(uint64(size))
	s.promTotalMsgs.Inc()
	s.promTotalBytes.Add(float64(size))
}

func (s *Stats) IncErr() {
	s.totalErr.Inc()
}
