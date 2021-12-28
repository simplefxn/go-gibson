package config

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Shopify/sarama"
	"github.com/simplefxn/kafkaLoad/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Metrics struct {
	Enable bool
	Port   int
	Path   string
}

type Service struct {
	Metrics        Metrics
	ReportInterval time.Duration
	Sarama         *sarama.Config
	Others         *SaramaComplex
}

var globalConf Service

func Get() *Service {
	return &globalConf
}

func init() {
	globalConf = Service{
		Sarama: sarama.NewConfig(),
		Others: &SaramaComplex{},
	}
}

const EnvPrefix = "KAFKALOADER"

func ParseCompression(scheme string) sarama.CompressionCodec {
	switch scheme {
	case "none":
		return sarama.CompressionNone
	case "gzip":
		return sarama.CompressionGZIP
	case "snappy":
		return sarama.CompressionSnappy
	case "lz4":
		return sarama.CompressionLZ4
	default:
		fmt.Fprintf(os.Stderr, "Unknown -compression: %s\n", scheme)
		os.Exit(1)
	}
	return sarama.CompressionNone
}

func ParsePartitioner(scheme string) sarama.PartitionerConstructor {
	scheme = strings.ToLower(scheme)
	switch scheme {
	case "manual":
		return sarama.NewManualPartitioner
	case "hash":
		return sarama.NewHashPartitioner
	case "random":
		return sarama.NewRandomPartitioner
	case "roundrobin":
		return sarama.NewRoundRobinPartitioner
	default:
		fmt.Fprintf(os.Stderr, "Unknown -partitioning: %s\n", scheme)
		os.Exit(1)
	}
	return sarama.NewHashPartitioner
}

func ParseIsolation(scheme string) sarama.IsolationLevel {
	scheme = strings.ToLower(scheme)
	switch scheme {
	case "readuncommitted":
		return sarama.ReadUncommitted
	case "readcommitted":
		return sarama.ReadCommitted
	default:
		fmt.Fprintf(os.Stderr, "Unknown -isolation: %s\n", scheme)
		os.Exit(1)
	}
	return sarama.ReadUncommitted
}

func ParseBalanceStrategy(scheme string) sarama.BalanceStrategy {
	scheme = strings.ToLower(scheme)
	switch scheme {
	case "range":
		return sarama.BalanceStrategyRange
	case "sticky":
		return sarama.BalanceStrategySticky
	case "roundrobin":
		return sarama.BalanceStrategyRoundRobin
	default:
		fmt.Fprintf(os.Stderr, "Unknown -rebalancestrategy: %s\n", scheme)
		os.Exit(1)
	}
	return sarama.BalanceStrategyRange
}

func ParseVersion(version string) *sarama.KafkaVersion {
	result, err := sarama.ParseKafkaVersion(version)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unknown -version: %s\n", version)
		os.Exit(1)
	}
	return &result
}

func ParseRequiredAcks(acks string) sarama.RequiredAcks {
	scheme := strings.ToLower(acks)
	switch scheme {
	case "noresponse":
		return sarama.NoResponse
	case "waitforlocal":
		return sarama.WaitForLocal
	case "waitforall":
		return sarama.WaitForAll
	}
	return sarama.NoResponse
}

func ParseOffsetsInitials(acks string) int64 {
	scheme := strings.ToLower(acks)
	switch scheme {
	case "offsetoldest":
		return sarama.OffsetOldest
	case "offsetnewest":
		return sarama.OffsetNewest
	}
	return sarama.OffsetNewest
}

func CreateTlsConfiguration(conf *Service) (t *tls.Config) {
	if conf.Others.Net.Cert != "" && conf.Others.Net.Key != "" && conf.Others.Net.CA != "" {
		cert, err := tls.LoadX509KeyPair(conf.Others.Net.Cert, conf.Others.Net.Key)
		if err != nil {
			logger.Log.Fatal(err)
		}

		caCert, err := os.ReadFile(conf.Others.Net.CA)
		if err != nil {
			logger.Log.Fatal(err)
		}

		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		t = &tls.Config{
			Certificates:       []tls.Certificate{cert},
			RootCAs:            caCertPool,
			InsecureSkipVerify: conf.Others.Net.VerifySSL,
		}
	}
	// will be nil by default if nothing is provided
	return t
}

type FakeData struct {
	UUID string `fake:"{uuid}" json:"uuid"`
	Ipv4 string `fake:"{ipv4address}" json:"ipv4"`
}

type SaramaComplex struct {
	Net struct {
		Host      string
		Port      int
		Cert      string
		Key       string
		CA        string
		VerifySSL bool
	}
	Producer struct {
		// The type of compression to use on messages (defaults to no compression).
		// Values  CompressionNone CompressionCodec CompressionGZIP CompressionSnappy CompressionLZ4 CompressionZSTD
		Compression string
		// Generates partitioners for choosing the partition to send messages to
		// (defaults to hashing the message key). Similar to the `partitioner.class`
		// setting for the JVM producer.
		// Values RandomPartitioner HashPartitioner RoundRobinPartitioner
		Partitioner string
		// RequiredAcks
		RequiredAcks string
		Generator    struct {
			Step     time.Duration
			Duration time.Duration
			Topic    string
			Verbose  bool
		}
	}

	// Consumer is the namespace for configuration related to consuming messages,
	// used by the Consumer.
	Consumer struct {
		Topic   string
		Verbose bool
		// Group is the namespace for configuring consumer group.
		Group struct {
			Name      string
			Rebalance struct {
				// Strategy for allocating topic partitions to members (default BalanceStrategyRange)
				// Values BalanceStrategySticky BalanceStrategyRoundRobin BalanceStrategyRange
				Strategy string
			}
		}
		Offsets struct {
			Initial string
		}
		// IsolationLevel support 2 mode:
		// 	- use `ReadUncommitted` (default) to consume and return all messages in message channel
		//	- use `ReadCommitted` to hide messages that are part of an aborted transaction
		// values ReadUncommitted ReadCommitted
		IsolationLevel string
	}

	Version string
}

// Bind each cobra flag to its associated viper configuration (config file and environment variable)
func BindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", EnvPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			fmt.Printf("%v:", f.Name)
			fmt.Printf("%v\n", val)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
