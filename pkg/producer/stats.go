package producer

import (
	"github.com/Shopify/sarama"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/atomic"
)

type StatsInterceptor struct {
	total            atomic.Uint64
	promTotalCounter prometheus.Counter
}

func newStats() *StatsInterceptor {
	// Initialize metrics counter
	s := &StatsInterceptor{
		promTotalCounter: promauto.NewCounter(prometheus.CounterOpts{
			Name: "gibson_producer_events_total",
			Help: "The total number produced events",
		}),
	}
	s.total.Store(0)
	return s
}

func (s *StatsInterceptor) GetTotal() uint64 {
	return s.total.Load()
}

func (s *StatsInterceptor) OnSend(msg *sarama.ProducerMessage) {
	s.total.Inc()
	s.promTotalCounter.Inc()
}
