package service

import (
	"context"
	"log"
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
	orderRepo     *storage.OrderRepository
}

type ServiceConfig struct {
	Cache         *storage.Cache
	KafkaConsumer *kafka.Consumer
	OrderRepo     *storage.OrderRepository
}

func NewService(cfg *ServiceConfig) *Service {
	return &Service{
		cache:         cfg.Cache,
		kafkaConsumer: cfg.KafkaConsumer,
		orderRepo:     cfg.OrderRepo,
	}
}

func (s *Service) StartConsumer() {
	log.Println("start consumer")
	go s.kafkaConsumer.Start(s.orderRepo)
}

func (s *Service) StopConsumer() error {
	log.Println("stop consumer")
	return s.kafkaConsumer.Stop()
}

func (s *Service) GetOrder(orderId string) (*storage.Order, error) {
	if order, ok := s.cache.Get(orderId); ok {
		return order, nil
	}

	order, err := s.orderRepo.GetByUID(context.Background(), orderId)
	if err != nil {
		return nil, err
	}

	s.cache.Set(orderId, *order)

	return order, nil
}
