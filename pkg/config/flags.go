package config

import (
	"time"

	"github.com/spf13/pflag"
)

func GetGenericFlags() *pflag.FlagSet {

	defaults := pflag.NewFlagSet("defaults for all commands", pflag.ExitOnError)
	defaults.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.kafkaLoad.yaml)")

	defaults.DurationVar(&globalConf.ReportInterval, "interval", 10*time.Second, "monitoring interval(60 seconds)")

	defaults.StringVar(&globalConf.Others.Net.Host, "kafka.broker", "cluster-kafka.kafka.svc", "host to connect or bind the socket")

	// Kafka.port Port use to connect to the kafka broker
	defaults.IntVar(&globalConf.Others.Net.Port, "kafka.port", 9092, "Port used to connect to the broker service")

	// Kafka.Cert is the cert to use if TLS is enabled
	defaults.StringVar(&globalConf.Others.Net.Cert, "kafka.cert", "", "server certificate to use for kafka connections, requires grpc_key, enables TLS")

	// Kafka.key is the key to use if TLS is enabled
	defaults.StringVar(&globalConf.Others.Net.Key, "kafka.key", "", "server private key to use for kafka connections, requires grpc_cert, enables TLS")

	// Kafka.ca	 is the CA to use if TLS is enabled
	defaults.StringVar(&globalConf.Others.Net.CA, "kafka.ca", "", "server CA to use for kafka connections, requires TLS, and enforces client certificate check")

	// Kafka.verifyssl Optional verify ssl certificates chain
	defaults.BoolVar(&globalConf.Others.Net.VerifySSL, "kafka.verifyssl", true, "Optional verify ssl certificates chain")

	// Enable prometheus stats
	defaults.BoolVar(&globalConf.Metrics.Enable, "metrics.enable", true, "enable prometheus metrics")

	// Prometheus metrics port to listen to
	defaults.IntVar(&globalConf.Metrics.Port, "metrics.port", 8080, "Port to liste for prometheus metrics scraping")
	defaults.StringVar(&globalConf.Metrics.Path, "metrics.path", "/metrics", "Path for prometheus metrics scraping ")

	return defaults
}

func GetProducerFlags() *pflag.FlagSet {
	producerFlags := pflag.NewFlagSet("producer flags", pflag.ExitOnError)

	// Producer
	producerFlags.IntVar(&globalConf.Sarama.Producer.MaxMessageBytes, "kafka.producer.maxmessagebytes", 1000000, "The maximum permitted size of a message (defaults to 1000000)")
	producerFlags.StringVar(&globalConf.Others.Producer.RequiredAcks, "kafka.producer.requiredacks", "noresponse", "The required number of acks needed from the broker [noresponse waitforlocal waitforall]")
	producerFlags.DurationVar(&globalConf.Sarama.Producer.Timeout, "kafka.producer.timeout", 10*time.Second, "The maximum duration the broker will wait the receipt of the number of RequiredAcks (defaults to 10 seconds)")
	producerFlags.StringVar(&globalConf.Others.Producer.Compression, "kafka.producer.compression", "none", "The type of compression to use on messages. [none gzip snappy lz4 zstd]")
	producerFlags.IntVar(&globalConf.Sarama.Producer.CompressionLevel, "kafka.producer.compressionlevel", 0, "The level of compression to use on messages")
	producerFlags.StringVar(&globalConf.Others.Producer.Partitioner, "kafka.producer.partitioner", "hash", "Generates partitioners for choosing the partition to send messages to. [random hash roundrobin]")
	producerFlags.BoolVar(&globalConf.Sarama.Producer.Idempotent, "kafka.producer.idempotent", false, "If enabled, the producer will ensure that exactly one copy of each message is written")
	producerFlags.BoolVar(&globalConf.Sarama.Producer.Return.Successes, "kafka.producer.return.successes", false, "If enabled, successfully delivered messages will be returned on the Successes channel")
	producerFlags.BoolVar(&globalConf.Sarama.Producer.Return.Errors, "kafka.producer.return.errors", true, "If enabled, messages that failed to deliver will be returned on the Errors channel")
	// Flush
	producerFlags.IntVar(&globalConf.Sarama.Producer.Flush.Bytes, "kafka.producer.flush.bytes", 0, "The best-effort number of bytes needed to trigger a flush")
	producerFlags.IntVar(&globalConf.Sarama.Producer.Flush.Messages, "kafka.producer.flush.messages", 0, "The best-effort number of messages needed to trigger a flush")
	producerFlags.DurationVar(&globalConf.Sarama.Producer.Flush.Frequency, "kafka.producer.flush.frequency", 500*time.Millisecond, "The best-effort frequency of flushes")
	producerFlags.IntVar(&globalConf.Sarama.Producer.Flush.MaxMessages, "kafka.producer.flush.maxmessages", 0, "The maximum number of messages the producer will send in a single broker request. Defaults to 0 for unlimited")
	// Retry
	producerFlags.IntVar(&globalConf.Sarama.Producer.Retry.Max, "kafka.producer.retry.max", 3, "The total number of times to retry sending a message (default 3)")
	producerFlags.DurationVar(&globalConf.Sarama.Producer.Retry.Backoff, "kafka.producer.retry.backoff", 100*time.Millisecond, "How long to wait for the cluster to settle between retries (default 100ms)")

	producerFlags.DurationVar(&globalConf.Others.Producer.Generator.Step, "kafka.producer.generator.step", 1*time.Second, "The time to take between sending messages in milliseconds (Default 1s)")
	producerFlags.DurationVar(&globalConf.Others.Producer.Generator.Duration, "kafka.producer.generator.duration", 60*time.Second, "Overall duration of the test. (Default 60s)")
	producerFlags.StringVar(&globalConf.Others.Producer.Generator.Topic, "kafka.producer.generator.topic", "topic1", "Topic to publish to kafka")
	producerFlags.BoolVar(&globalConf.Others.Producer.Generator.Verbose, "kafka.producer.generator.verbose", false, "Output messages to the log")

	return producerFlags
}

func GetConsumerFlags() *pflag.FlagSet {
	consumerFlags := pflag.NewFlagSet("consumer flags", pflag.ExitOnError)
	// Consumer
	consumerFlags.StringVar(&globalConf.Others.Consumer.Group.Name, "kafka.consumer.group.name", "my-group-1", "the consumer group for a consumer")
	consumerFlags.StringVar(&globalConf.Others.Consumer.Topic, "kafka.consumer.topic", "my-group-1", "the consumer group for a consumer")
	consumerFlags.BoolVar(&globalConf.Others.Consumer.Verbose, "kafka.consumer.verbose", false, "output consumed messages in the log")

	consumerFlags.DurationVar(&globalConf.Sarama.Consumer.Group.Session.Timeout, "kafka.consumer.group.session.timeout", 10*time.Second, "The timeout used to detect consumer failures when using Kafka's group management facility")

	consumerFlags.DurationVar(&globalConf.Sarama.Consumer.Group.Heartbeat.Interval, "kafka.consumer.group.heartbeat.interval", 3*time.Second, "The expected time between heartbeats to the consumer coordinator when using Kafka's group management facilities")
	// Rebalance
	consumerFlags.StringVar(&globalConf.Others.Consumer.Group.Rebalance.Strategy, "kafka.consumer.group.rebalance.strategy", "range", "Strategy for allocating topic partitions to members. Values [Range Sticky RoundRobin](Default to Rage)")
	consumerFlags.DurationVar(&globalConf.Sarama.Consumer.Group.Rebalance.Timeout, "kafka.consumer.group.rebalance.timeout", 60*time.Second, "The maximum allowed time for each worker to join the group once a rebalance has begun")

	consumerFlags.IntVar(&globalConf.Sarama.Consumer.Group.Rebalance.Retry.Max, "kafka.consumer.group.rebalance.retry.max", 4, "When a new consumer joins a consumer group the set of consumers attempt to rebalance the load to assign partitions to each consumer")
	consumerFlags.DurationVar(&globalConf.Sarama.Consumer.Group.Rebalance.Retry.Backoff, "kafka.consumer.group.rebalance.retry.backoff", 2*time.Second, "Backoff time between retries during rebalance (default 2s)")
	// Retry
	consumerFlags.DurationVar(&globalConf.Sarama.Consumer.Retry.Backoff, "kafka.consumer.retry.backoff", 2*time.Second, "How long to wait after a failing to read from a partition before trying again (default 2s)")
	// Fetch
	consumerFlags.Int32Var(&globalConf.Sarama.Consumer.Fetch.Min, "kafka.consumer.fetch.min", 1, "The minimum number of message bytes to fetch in a request")
	consumerFlags.Int32Var(&globalConf.Sarama.Consumer.Fetch.Default, "kafka.consumer.fetch.default", 1073741824, "The default number of message bytes to fetch from the broker in each request (default 1MB)")
	consumerFlags.Int32Var(&globalConf.Sarama.Consumer.Fetch.Max, "kafka.consumer.fetch.max", 0, "The maximum number of message bytes to fetch from the broker in a single request. Defaults to 0 (no limit)")

	consumerFlags.DurationVar(&globalConf.Sarama.Consumer.MaxWaitTime, "kafka.consumer.maxwaittime", 250*time.Millisecond, "The maximum amount of time the broker will wait for Consumer.Fetch.Min bytes to become available before it returns fewer than that anyways")
	consumerFlags.DurationVar(&globalConf.Sarama.Consumer.MaxProcessingTime, "kafka.consumer.maxprocessingtime", 100*time.Millisecond, "The maximum amount of time the consumer expects a message takes to process for the user")
	// Return
	consumerFlags.BoolVar(&globalConf.Sarama.Consumer.Return.Errors, "kafka.consumer.return.errors", false, "If enabled, any errors that occurred while consuming are returned on the Errors channel (default disabled)")
	// Offsets
	consumerFlags.BoolVar(&globalConf.Sarama.Consumer.Offsets.AutoCommit.Enable, "kafka.consumer.offsets.autocommit.enable", true, "Whether or not to auto-commit updated offsets back to the broker. (default enabled)")
	consumerFlags.DurationVar(&globalConf.Sarama.Consumer.Offsets.AutoCommit.Interval, "kafka.consumer.offsets.autocommit.interval", 1*time.Second, "How frequently to commit updated offsets. Ineffective unless auto-commit is enabled")
	// TODO , needs value of OffsetNewest or OffsetOldest in int64 ??
	consumerFlags.StringVar(&globalConf.Others.Consumer.Offsets.Initial, "kafka.consumer.offsets.initial", "offsetnewest", "The initial offset to use if no offset was previously committed.Should be offsetnewest or offsetoldest. Defaults to OffsetNewest")
	consumerFlags.DurationVar(&globalConf.Sarama.Consumer.Offsets.Retention, "kafka.consumer.offsets.retention", 0*time.Second, "The retention duration for committed offsets. If zero, disabled (in which case the `offsets.retention.minutes` option on the broker will be used)")

	consumerFlags.IntVar(&globalConf.Sarama.Consumer.Offsets.Retry.Max, "kafka.consumer.offsets.retry.max", 3, "The total number of times to retry failing commit requests during OffsetManager shutdown (default 3)")

	consumerFlags.StringVar(&globalConf.Others.Consumer.IsolationLevel, "kafka.consumer.isolationlevel", "readuncommitted", "Isolation Level (Default readUncommitted) Values [readuncommitted readcommitted]")

	return consumerFlags
}
