package producer

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
	"github.com/simplefxn/go-gibson/pkg/metrics"
	"github.com/spf13/cobra"
)

type Gibson struct {
	producer sarama.AsyncProducer
}

func New(cmd *cobra.Command, conf *config.Service) (*Gibson, error) {
	config.SetLogLevel(cmd)

	producer := Gibson{
		producer: newProducer(conf),
	}

	compression, err := cmd.Flags().GetString("kafka.producer.compression")
	if err != nil {
		return nil, err
	}
	config.Get().Sarama.Producer.Compression = config.ParseCompression(compression)

	partitioner, err := cmd.Flags().GetString("kafka.producer.partitioner")
	if err != nil {
		return nil, err
	}

	config.Get().Sarama.Producer.Partitioner = config.ParsePartitioner(partitioner)

	version, err := cmd.Flags().GetString("kafka.version")
	if err != nil {
		return nil, err
	}

	config.Get().Sarama.Version = *config.ParseVersion(version)

	// Initialize metrics counter
	metrics.ProducerEventCounter.Store(0)

	config.Dump(cmd)

	return &producer, nil
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
