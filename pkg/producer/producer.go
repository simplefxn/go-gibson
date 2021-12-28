package producer

import (
	"fmt"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
	"github.com/spf13/cobra"
)

type Gibson struct {
	producer sarama.AsyncProducer
	stats    *StatsInterceptor
}

func New(cmd *cobra.Command, conf *config.Service) (*Gibson, error) {
	config.SetLogLevel(cmd)

	stats := newStats()

	producer := Gibson{
		producer: newProducer(conf, stats),
		stats:    stats,
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

	config.Dump(cmd)

	return &producer, nil
}

func (r Gibson) Input() chan<- *sarama.ProducerMessage {
	return r.producer.Input()
}

func (r Gibson) Close() error {
	logger.Log.Infof("Total messages: %v", r.stats.GetTotal())
	return r.producer.Close()
}

func newProducer(conf *config.Service, stats *StatsInterceptor) sarama.AsyncProducer {
	sarama.Logger = logger.NewSaramaLogger(logger.GetLogger())

	tlsConfig := config.CreateTlsConfiguration(conf)
	if tlsConfig != nil {
		config.Get().Sarama.Net.TLS.Enable = true
		config.Get().Sarama.Net.TLS.Config = tlsConfig
	}

	brokers := []string{fmt.Sprintf("%s:%d", conf.Others.Net.Host, conf.Others.Net.Port)}
	interceptors := []sarama.ProducerInterceptor{
		stats,
	}
	if strings.ToLower(config.Get().Globals.LogLevel) == "debug" {
		interceptors = append(interceptors, newLog())
	}
	config.Get().Sarama.Producer.Interceptors = interceptors
	producer, err := sarama.NewAsyncProducer(brokers, config.Get().Sarama)

	if err != nil {
		logger.Log.Fatalf("Failed to start Sarama producer: %v", err)
	}
	logger.Log.Infof("Created async producer")
	return producer
}
