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

type MessageHandler interface {
	HandleMessage(value []byte, offset kafka.Offset, repo *storage.OrderRepository) error
}

type Consumer struct {
	consumer *kafka.Consumer
	handler  MessageHandler
	stop     bool
}

func NewConsumer(h MessageHandler, conf config.Config) (*Consumer, error) {

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

func (c *Consumer) Start(repo *storage.OrderRepository) {
	c.stop = false
	go func() {
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

			if err = c.handler.HandleMessage(kafkaMsg.Value, kafkaMsg.TopicPartition.Offset, repo); err != nil {
				log.Println(err)
				continue
			}

			if _, err = c.consumer.StoreMessage(kafkaMsg); err != nil {
				log.Println(err)
				continue
			}
		}
	}()
}

func (c *Consumer) Stop() error {
	c.stop = true
	if _, err := c.consumer.Commit(); err != nil {
		return err
	}
	return c.consumer.Close()
}
