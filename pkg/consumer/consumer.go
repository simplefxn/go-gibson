package consumer

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
	"github.com/spf13/cobra"
)

// Consumer represents a Sarama consumer group consumer
type Gibson struct {
	ready    chan bool
	client   sarama.ConsumerGroup
	topic    string
	ctx      context.Context
	callback func(msg string) error
	stats    *StatsInterceptor
}

func New(ctx context.Context, cmd *cobra.Command, conf *config.Service, callback func(msg string) error) (*Gibson, error) {

	config.SetLogLevel(cmd)

	sarama.Logger = logger.NewSaramaLogger(logger.GetLogger())

	rebalance, err := cmd.Flags().GetString("kafka.consumer.group.rebalance.strategy")
	if err != nil {
		return nil, err
	}

	config.Get().Sarama.Consumer.Group.Rebalance.Strategy = config.ParseBalanceStrategy(rebalance)

	isolation, err := cmd.Flags().GetString("kafka.consumer.isolationlevel")
	if err != nil {
		return nil, err
	}

	config.Get().Sarama.Consumer.IsolationLevel = config.ParseIsolation(isolation)

	version, err := cmd.Flags().GetString("kafka.version")
	if err != nil {
		return nil, err
	}

	config.Get().Sarama.Version = *config.ParseVersion(version)

	offsetInitial, err := cmd.Flags().GetString("kafka.consumer.offsets.initial")
	if err != nil {
		return nil, err
	}

	config.Get().Sarama.Consumer.Offsets.Initial = config.ParseOffsetsInitials(offsetInitial)

	config.Get().Sarama.Version = *config.ParseVersion(version)

	tlsConfig := config.CreateTlsConfiguration(conf)
	if tlsConfig != nil {
		config.Get().Sarama.Net.TLS.Enable = true
		config.Get().Sarama.Net.TLS.Config = tlsConfig
	}

	brokers := []string{fmt.Sprintf("%s:%d", conf.Others.Net.Host, conf.Others.Net.Port)}

	stats := newStats()

	// Configure interceptors
	interceptors := []sarama.ConsumerInterceptor{
		stats,
	}
	if strings.ToLower(config.Get().Globals.LogLevel) == "debug" {
		interceptors = append(interceptors, newLog())
	}
	config.Get().Sarama.Consumer.Interceptors = interceptors

	client, err := sarama.NewConsumerGroup(brokers, config.Get().Others.Consumer.Group.Name, config.Get().Sarama)
	if err != nil {
		return nil, fmt.Errorf("error creating consumer group client: %v", err)
	}

	consumer := Gibson{
		ready:    make(chan bool),
		client:   client,
		topic:    config.Get().Others.Consumer.Topic,
		ctx:      ctx,
		callback: callback,
		stats:    stats,
	}

	config.Dump(cmd)

	return &consumer, nil
}

func (c *Gibson) Run() {
	var wg sync.WaitGroup
	var err error

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			// `Consume` should be called inside an infinite loop, when a
			// server-side rebalance happens, the consumer session will need to be
			// recreated to get the new claims
			if err := c.client.Consume(c.ctx, []string{c.topic}, c); err != nil {
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if c.ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready // Await till the consumer has been set up
	log.Println("Sarama consumer up and running!...")

	<-c.ctx.Done()
	log.Println("terminating: context cancelled")

	wg.Wait()

	logger.Log.Infof("Total messages: %v", c.stats.GetTotal())

	if err = c.client.Close(); err != nil {
		logger.Log.Errorf("Error closing client: %v", err)
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Gibson) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Gibson) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Gibson) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		if err := c.callback(string(message.Value)); err != nil {
			logger.Log.Info(err)
		}
		session.MarkMessage(message, "")
	}

	return nil
}
