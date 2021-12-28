package producer

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
)

type Gibson struct {
	producer sarama.AsyncProducer
}

func New(conf *config.Service) *Gibson {
	risProducer := Gibson{
		producer: newProducer(conf),
	}

	return &risProducer
}

func (r Gibson) Input() chan<- *sarama.ProducerMessage {
	return r.producer.Input()
}

func (r Gibson) Close() error {
	return r.producer.Close()
}

func newProducer(conf *config.Service) sarama.AsyncProducer {
	sarama.Logger = logger.NewSaramaLogger(logger.GetLogger())

	tlsConfig := config.CreateTlsConfiguration(conf)
	if tlsConfig != nil {
		config.Get().Sarama.Net.TLS.Enable = true
		config.Get().Sarama.Net.TLS.Config = tlsConfig
	}

	brokers := []string{fmt.Sprintf("%s:%d", conf.Others.Net.Host, conf.Others.Net.Port)}
	producer, err := sarama.NewAsyncProducer(brokers, config.Get().Sarama)

	if err != nil {
		logger.Log.Fatalf("Failed to start Sarama producer: %v", err)
	}
	logger.Log.Infof("Created async producer")
	return producer
}
