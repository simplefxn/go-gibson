package consumer

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
)

// Consumer represents a Sarama consumer group consumer
type Consumer struct {
	ready    chan bool
	client   sarama.ConsumerGroup
	topic    string
	ctx      context.Context
	callback func(msg string) error
}

func New(ctx context.Context, conf *config.Service, callback func(msg string) error) (*Consumer, error) {

	sarama.Logger = logger.NewSaramaLogger(logger.GetLogger())

	tlsConfig := config.CreateTlsConfiguration(conf)
	if tlsConfig != nil {
		config.Get().Sarama.Net.TLS.Enable = true
		config.Get().Sarama.Net.TLS.Config = tlsConfig
	}

	brokers := []string{fmt.Sprintf("%s:%d", conf.Others.Net.Host, conf.Others.Net.Port)}
	client, err := sarama.NewConsumerGroup(brokers, config.Get().Others.Consumer.Group.Name, config.Get().Sarama)
	if err != nil {
		return nil, fmt.Errorf("error creating consumer group client: %v", err)
	}

	consumer := Consumer{
		ready:    make(chan bool),
		client:   client,
		topic:    config.Get().Others.Consumer.Topic,
		ctx:      ctx,
		callback: callback,
	}

	return &consumer, nil
}

func (c *Consumer) Run() {
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
	if err = c.client.Close(); err != nil {
		logger.Log.Errorf("Error closing client: %v", err)
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/main/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		session.MarkMessage(message, "")
		//log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		if err := c.callback(string(message.Value)); err != nil {
			logger.Log.Info(err)
		}

	}

	return nil
}
