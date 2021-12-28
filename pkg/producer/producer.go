package producer

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
)

type RIS struct {
	RisEventProducer sarama.AsyncProducer
}

func New(conf *config.Service) *RIS {
	risProducer := RIS{
		RisEventProducer: newProducer(conf),
	}

	return &risProducer
}

func (r RIS) Input() chan<- *sarama.ProducerMessage {
	return r.RisEventProducer.Input()
}

func (r RIS) Close() error {
	return r.RisEventProducer.Close()
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
