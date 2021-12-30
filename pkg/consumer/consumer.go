package consumer

import (
	"context"
	"fmt"
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

func New(ctx context.Context, cmd *cobra.Command, callback func(msg string) error) (*Gibson, error) {
	conf := config.Get()

	sarama.Logger = logger.NewSaramaLogger(logger.GetLogger())

	config.SetLogLevel(cmd)

	err := config.SetConsumerFlags(cmd.Flags())
	if err != nil {
		return nil, err
	}

	tlsConfig := config.CreateTlsConfiguration(conf)
	if tlsConfig != nil {
		config.Get().Sarama.Net.TLS.Enable = true
		config.Get().Sarama.Net.TLS.Config = tlsConfig
	}

	brokers := []string{fmt.Sprintf("%s:%d", conf.Others.Net.Host, conf.Others.Net.Port)}

	// Configure interceptors
	stats := newStats()

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
		topic:    config.GetConsumerTopic(),
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
				logger.Log.Fatalf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if c.ctx.Err() != nil {
				return
			}
			c.ready = make(chan bool)
		}
	}()

	<-c.ready // Await till the consumer has been set up
	logger.Log.Info("Sarama consumer up and running!...")

	<-c.ctx.Done()
	logger.Log.Info("terminating: context cancelled")

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
			return err
		}
		session.MarkMessage(message, "")
	}

	return nil
}
