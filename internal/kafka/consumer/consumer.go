package kafka

import (
	"log"
	"wb_tech_L0/internal/config"
	"wb_tech_L0/internal/storage"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	sessionTimeout = 7000 // ms
	noTimeout      = -1
)

type Handler interface {
	HandleMessage(message []byte, offset kafka.Offset, cache *storage.Cache) error
}

type Consumer struct {
	consumer *kafka.Consumer
	handler  Handler
	stop     bool
}

func NewConsumer(h Handler, conf config.Config) (*Consumer, error) {

	confMap := &kafka.ConfigMap{
		"bootstrap.servers":        conf.Kafka.KafkaAddress,
		"group.id":                 conf.Kafka.ConsumerGroup,
		"session.timeout.ms":       sessionTimeout,
		"enable.auto.offset.store": false,
		"enable.auto.commit":       true,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest",
	}

	c, err := kafka.NewConsumer(confMap)
	if err != nil {
		return nil, err
	}

	if err = c.Subscribe(conf.Kafka.Topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		consumer: c,
		handler:  h,
	}, nil
}

func (c *Consumer) Start(cache *storage.Cache) {
	c.stop = false
	for {
		if c.stop {
			break
		}
		kafkaMsg, err := c.consumer.ReadMessage(noTimeout)
		if err != nil {
			log.Println(err)
		}

		if kafkaMsg == nil {
			continue
		}

		if err = c.handler.HandleMessage(kafkaMsg.Value, kafkaMsg.TopicPartition.Offset, cache); err != nil {
			log.Println(err)
			continue
		}

		if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
			log.Println(err)
			continue
		}
	}
}

func (c *Consumer) Stop() error {
	c.stop = true
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}
	return c.consumer.Close()
}
