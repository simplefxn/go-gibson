package consumer

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

func (s *LogInterceptor) OnConsume(msg *sarama.ConsumerMessage) {
	logger.Log.Infof("%s: %s", msg.Topic, string(msg.Value))
}
