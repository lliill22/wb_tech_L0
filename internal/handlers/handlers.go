package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"wb_tech_L0/internal/storage"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaHandler interface {
	HandleMessage(message []byte, offset kafka.Offset, repo *storage.OrderRepository) error
}

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleMessage(message []byte, offset kafka.Offset, repo *storage.OrderRepository) error {
	order := &storage.Order{}
	if err := json.Unmarshal(message, order); err != nil {
		return err
	}

	ctx := context.Background()

	repo.Insert(ctx, order)
	log.Panicln("Massege was write in repository")

	return nil
}

func (r *Router) GetOrder(w http.ResponseWriter, req *http.Request) {

}
