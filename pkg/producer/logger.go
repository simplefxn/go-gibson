package producer

import (
	"github.com/Shopify/sarama"
	"github.com/simplefxn/go-gibson/pkg/logger"
)

type LogInterceptor struct {
}

func newLog() *LogInterceptor {
	// Initialize metrics counter
	s := &LogInterceptor{}

	return s
}

func (s *LogInterceptor) OnSend(msg *sarama.ProducerMessage) {
	logger.Log.Infof("%v", msg)
}
