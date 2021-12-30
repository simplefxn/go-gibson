package admin

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/Shopify/sarama"
	"github.com/simplefxn/go-gibson/pkg/config"
	"github.com/simplefxn/go-gibson/pkg/logger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type ClusterAdmin struct {
	admin sarama.ClusterAdmin
	conf  *config.Service
	data  *Data
}
type Topic struct {
	Name              string `yaml:"Name"`
	Partitions        int    `yaml:"Partitions"`
	ReplicationFactor int    `yaml:"Replication_factor"`
}

type Data struct {
	Kafka struct {
		Topics []Topic `yaml:"Topics"`
	} `yaml:"kafka"`
}

func New(ctx context.Context, cmd *cobra.Command) (*ClusterAdmin, error) {
	config.SetLogLevel(cmd)

	conf := config.Get()

	sarama.Logger = logger.NewSaramaLogger(logger.GetLogger())

	tlsConfig := config.CreateTlsConfiguration(conf)
	if tlsConfig != nil {
		config.Get().Sarama.Net.TLS.Enable = true
		config.Get().Sarama.Net.TLS.Config = tlsConfig
	}

	brokers := []string{fmt.Sprintf("%s:%d", conf.Others.Net.Host, conf.Others.Net.Port)}

	admin, err := sarama.NewClusterAdmin(brokers, config.Get().Sarama)
	if err != nil {
		return nil, err
	}

	yamlFile, err := ioutil.ReadFile(viper.GetViper().ConfigFileUsed())
	if err != nil {
		return nil, fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	d := &Data{}
	err = yaml.Unmarshal(yamlFile, d)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %v", err)
	}

	ca := ClusterAdmin{
		admin: admin,
		conf:  conf,
		data:  d,
	}

	return &ca, nil
}

func (c *ClusterAdmin) Process() error {

	if c.admin == nil {
		return fmt.Errorf("call Process with an invalid cluster admin")
	}

	for _, topic := range c.data.Kafka.Topics {
		logger.Log.Debugf("Creating topic %s", topic.Name)
		err := c.admin.CreateTopic(topic.Name, &sarama.TopicDetail{
			NumPartitions:     int32(topic.Partitions),
			ReplicationFactor: int16(topic.ReplicationFactor),
		}, false)
		if err != nil {
			logger.Log.Debug(err)
		}
	}
	c.admin.Close()

	return nil
}
