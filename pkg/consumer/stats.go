package consumer

import (
	"github.com/Shopify/sarama"
	"go.uber.org/atomic"
)

type StatsInterceptor struct {
	total atomic.Uint64
}

func newStats() *StatsInterceptor {
	// Initialize metrics counter
	s := &StatsInterceptor{}
	s.total.Store(0)
	return s
}

func (s *StatsInterceptor) GetTotal() uint64 {
	return s.total.Load()
}

func (s *StatsInterceptor) OnConsume(msg *sarama.ConsumerMessage) {
	s.total.Inc()
}
