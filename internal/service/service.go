package service

import (
	kafka "wb_tech_L0/internal/kafka/consumer"
	"wb_tech_L0/internal/storage"
)

type ServiceImpl interface {
	StartConsumer()
	StopConsumer() error
}

type Service struct {
	cache         *storage.Cache
	kafkaConsumer *kafka.Consumer
}

type ServiceConfig struct {
	Cache         *storage.Cache
	KafkaConsumer *kafka.Consumer
	DB            *storage.Database
}

func NewService(cfg *ServiceConfig) ServiceImpl {
	return &Service{
		cache:         cfg.Cache,
		kafkaConsumer: cfg.KafkaConsumer,
	}
}

func (s *Service) StartConsumer() {
	go s.kafkaConsumer.Start(s.cache)
}

func (s *Service) StopConsumer() error {
	return s.kafkaConsumer.Stop()
}
