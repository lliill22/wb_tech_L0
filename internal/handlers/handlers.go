package handlers

import (
	"encoding/json"
	"log"
	"wb_tech_L0/internal/storage"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaHandler interface {
	HandleMessage(message []byte, offset kafka.Offset) error
}

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleMessage(message []byte, offset kafka.Offset, cache *storage.Cache) error {
	order := &storage.Order{}
	if err := json.Unmarshal(message, order); err != nil {
		return err
	}
	cache.Set(order.OrderUID, *order)

	log.Panicln("Massege was write in cache")

	return nil
}
