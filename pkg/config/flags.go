package config

import (
	"time"

	"github.com/spf13/pflag"
)

func GetGenericFlags() *pflag.FlagSet {

	defaults := pflag.NewFlagSet("defaults for all commands", pflag.ExitOnError)
	defaults.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config.yaml)")

	defaults.StringVar(&globalConf.Globals.LogLevel, "log.level", "info", "verbosity of logs, known levels are: debug, info, warn, error, fatal, panic")

	return defaults
}

func GetKafkaGenericFlags() *pflag.FlagSet {

	defaults := pflag.NewFlagSet("defaults for all commands", pflag.ExitOnError)
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

	defaults.StringVar(&globalConf.Others.Version, "kafka.version", "3.0.0", "The version of Kafka that Sarama will assume it is running against. Defaults to the oldest supported stable version")

	defaults.StringVar(&globalConf.Sarama.ClientID, "kafka.clientid", "gibson-client", "A user-provided string sent with every request to the brokers for logging, debugging, and auditing purposes")
	defaults.StringVar(&globalConf.Sarama.RackID, "kafka.rackid", "rack01", "A rack identifier for this client.")

	defaults.IntVar(&globalConf.Sarama.ChannelBufferSize, "kafka.channelbuffersize", 1024, "The number of events to buffer in internal and external channels")
	defaults.BoolVar(&globalConf.Sarama.ApiVersionsRequest, "kafka.apiversionrequest", true, "The version of Kafka that Sarama will assume it is running against")

	return defaults
}

func GetMetricsFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("Metrics parameters", pflag.ExitOnError)
	// Enable prometheus stats
	flags.BoolVar(&globalConf.Metrics.Enable, "metrics.enable", true, "enable prometheus metrics")

	// Prometheus metrics port to listen to
	flags.IntVar(&globalConf.Metrics.Port, "metrics.port", 8080, "Port to liste for prometheus metrics scraping")
	flags.StringVar(&globalConf.Metrics.Path, "metrics.path", "/metrics", "Path for prometheus metrics scraping ")

	return flags
}

func GetSSEFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("SSE parameters", pflag.ExitOnError)

	flags.StringVar(&globalConf.SSE.Section.URL, "sse.url", "https://ris-live.ripe.net/v1/stream", "sse url to connect for sse")

	flags.StringVar(&globalConf.SSE.Section.ClientId, "sse.clientstring", "ris-client", "sse url to connect for sse")
	flags.DurationVar(&globalConf.SSE.Section.IdleTimeout, "sse.idletimeout", 10*time.Second, "sse idle timeout")
	flags.DurationVar(&globalConf.SSE.Section.ReportInterval, "sse.report.interval", 60*time.Second, "report to console interval(60 seconds)")

	return flags
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

	return producerFlags
}
func GetKafkaAdminFlags() *pflag.FlagSet {
	adminFlags := pflag.NewFlagSet("admin flags", pflag.ExitOnError)

	adminFlags.IntVar(&globalConf.Sarama.Admin.Retry.Max, "kafka.admin.retry.max", 5, "The total number of times to retry sending (retriable) admin requests (default 5)")
	adminFlags.DurationVar(&globalConf.Sarama.Admin.Retry.Backoff, "kafka.admin.retry.backoff", 100*time.Millisecond, "Backoff time between retries of a failed request (default 100ms)")
	adminFlags.DurationVar(&globalConf.Sarama.Admin.Timeout, "kafka.admin.timeout", 3*time.Second, "The maximum duration the administrative Kafka client will wait for ClusterAdmin operations, including topics, brokers, configurations and ACLs (defaults to 3 seconds)")

	return adminFlags
}

func GetKafkaNetFlags() *pflag.FlagSet {
	netFlags := pflag.NewFlagSet("net flags", pflag.ExitOnError)

	netFlags.IntVar(&globalConf.Sarama.Net.MaxOpenRequests, "kafka.net.maxopenrequests", 5, "How many outstanding requests a connection is allowed to have before sending on it blocks (default 5)")
	netFlags.DurationVar(&globalConf.Sarama.Net.DialTimeout, "kafka.net.dialtimeout", 30*time.Second, "How long to wait for the initial connection")
	netFlags.DurationVar(&globalConf.Sarama.Net.ReadTimeout, "kafka.net.readtimeout", 30*time.Second, "How long to wait for a response")
	netFlags.DurationVar(&globalConf.Sarama.Net.WriteTimeout, "kafka.net.writetimeout", 30*time.Second, "How long to wait for a transmit")

	return netFlags
}

func GetKafkaMetadataFlags() *pflag.FlagSet {
	metaFlags := pflag.NewFlagSet("metadata flags", pflag.ExitOnError)

	metaFlags.IntVar(&globalConf.Sarama.Metadata.Retry.Max, "kafka.metadata.retry.max", 3, "The total number of times to retry a metadata request when the cluster is in the middle of a leader election (default 3)")
	metaFlags.DurationVar(&globalConf.Sarama.Metadata.Retry.Backoff, "kafka.metadata.retry.backoff", 250*time.Microsecond, "How long to wait for leader election to occur before retrying (default 250ms)")
	metaFlags.DurationVar(&globalConf.Sarama.Metadata.RefreshFrequency, "kafka.metadata.refreshfrequency", 10*time.Minute, "How frequently to refresh the cluster metadata in the background. Defaults to 10 minutes")
	metaFlags.BoolVar(&globalConf.Sarama.Metadata.Full, "kafka.metadata.full", true, "Whether to maintain a full set of metadata for all topics, or just he minimal set that has been necessary so far.")

	return metaFlags
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

func SetConsumerFlags(flags *pflag.FlagSet) error {
	rebalance, err := flags.GetString("kafka.consumer.group.rebalance.strategy")
	if err != nil {
		return err
	}

	globalConf.Sarama.Consumer.Group.Rebalance.Strategy = parseBalanceStrategy(rebalance)

	isolation, err := flags.GetString("kafka.consumer.isolationlevel")
	if err != nil {
		return err
	}

	globalConf.Sarama.Consumer.IsolationLevel = parseIsolation(isolation)

	version, err := flags.GetString("kafka.version")
	if err != nil {
		return err
	}

	globalConf.Sarama.Version = *parseVersion(version)

	offsetInitial, err := flags.GetString("kafka.consumer.offsets.initial")
	if err != nil {
		return err
	}

	globalConf.Sarama.Consumer.Offsets.Initial = parseOffsetsInitials(offsetInitial)

	return nil
}

func SetProducerFlags(flags *pflag.FlagSet) error {
	compression, err := flags.GetString("kafka.producer.compression")
	if err != nil {
		return err
	}
	globalConf.Sarama.Producer.Compression = ParseCompression(compression)

	partitioner, err := flags.GetString("kafka.producer.partitioner")
	if err != nil {
		return err
	}

	globalConf.Sarama.Producer.Partitioner = ParsePartitioner(partitioner)

	version, err := flags.GetString("kafka.version")
	if err != nil {
		return err
	}

	globalConf.Sarama.Version = *parseVersion(version)

	return nil
}

func GetSenderFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("UDP Sender", pflag.ExitOnError)

	flags.StringVar(&globalConf.Sender.Address, "sender.address", "sender.default.cluster.local", "service to send the messages to")
	flags.IntVar(&globalConf.Sender.Port, "sender.port", 5000, "port to send the messages to")

	return flags
}

func GetReceiverFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("UDP Receiver", pflag.ExitOnError)

	flags.StringVar(&globalConf.Receiver.Address, "receiver.address", "receiver.default.cluster.local", "service who will receive the messages")
	flags.IntVar(&globalConf.Receiver.Port, "receiver.port", 5000, "port to receive messages")

	return flags
}

func GetNatsGenericFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("Nats Generic", pflag.ExitOnError)

	flags.StringVar(&globalConf.Nats.URL, "nats.server", "http://nats.svc.default.cluster.local:9090", "nats server url")
	flags.StringVar(&globalConf.Nats.CA, "nats.ca", "", "ca pem file for nats server")
	// Kafka.Cert is the cert to use if TLS is enabled
	flags.StringVar(&globalConf.Nats.Cert, "nats.cert", "", "server certificate to use for nats connections")

	// Kafka.key is the key to use if TLS is enabled
	flags.StringVar(&globalConf.Nats.Key, "nats.key", "", "server private key to use for nats connections")

	flags.BoolVar(&globalConf.Nats.VerifySSL, "nats.verifyssl", true, "Optional verify ssl certificates chain")
	flags.DurationVar(&globalConf.Nats.ReportInterval, "nats.report.interval", 60*time.Second, "report to console interval(60 seconds)")

	return flags
}

func GetNatsPublisherFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("Nats Publisher", pflag.ExitOnError)

	flags.StringVar(&globalConf.Nats.Publisher.Topic, "nats.publisher.topic", "topic1", "publishing topic to nats server")
	return flags
}

func GetNatsSubscriberFlags() *pflag.FlagSet {
	flags := pflag.NewFlagSet("Nats Subscriber", pflag.ExitOnError)

	flags.StringVar(&globalConf.Nats.Subscriber.Topic, "nats.subscriber.topic", "topic1", "subscriber topic to nats server")
	return flags
}
