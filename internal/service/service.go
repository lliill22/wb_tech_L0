package service

import (
	"context"
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

func NewService(cfg *ServiceConfig) ServiceImpl {
	return &Service{
		cache:         cfg.Cache,
		kafkaConsumer: cfg.KafkaConsumer,
		orderRepo:     cfg.OrderRepo,
	}
}

func (s *Service) StartConsumer() {
	go s.kafkaConsumer.Start(s.orderRepo)
}

func (s *Service) StopConsumer() error {
	return s.kafkaConsumer.Stop()
}

func (s *Service) GetOrder(orderId string) (storage.Order, error) {
	if order, ok := s.cache.Get(orderId); ok {
		return order, nil
	}

	ctx := context.Background()
	order, err := s.orderRepo.GetByUID(ctx, orderId)
	if err != nil {
		return storage.Order{}, err
	}

	s.cache.Set(orderId, *order)

	return *order, nil
}
