package admin

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/Shopify/sarama"
	"github.com/simplefxn/kafkaLoad/pkg/config"
	"github.com/simplefxn/kafkaLoad/pkg/logger"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type ClusterAdmin struct {
	admin sarama.ClusterAdmin
	conf  *config.Service
	data  *Data
}

type Data struct {
	Kafka struct {
		Topic struct {
			Name              string `yaml:"Name"`
			Partitions        int    `yaml:"Partitions"`
			ReplicationFactor int    `yaml:"Replication_factor"`
		} `yaml:"Topic"`
	} `yaml:"kafka"`
}

func New(ctx context.Context, conf *config.Service) (*ClusterAdmin, error) {

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
		logger.Log.Info(err)
		return nil, fmt.Errorf("yamlFile.Get err   #%v ", err)
	}
	d := &Data{}
	err = yaml.Unmarshal(yamlFile, d)
	if err != nil {
		logger.Log.Info(err)
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
	topic := c.data.Kafka.Topic

	logger.Log.Infof("Creating topic %s", topic.Name)
	err := c.admin.CreateTopic(topic.Name, &sarama.TopicDetail{
		NumPartitions:     int32(topic.Partitions),
		ReplicationFactor: int16(topic.ReplicationFactor),
	}, false)
	if err != nil {
		return err
	}

	logger.Log.Info("Letting topics settle in (10s)")
	time.Sleep(10 * time.Second)
	c.admin.Close()

	return nil
}
