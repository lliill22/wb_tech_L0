package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"wb_tech_L0/internal/storage"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/go-chi/chi"
)

type Handler struct {
}

func (h *Handler) HandleMessage(value []byte, offset kafka.Offset, repo *storage.OrderRepository) error {
	log.Println("start handle message")
	var order storage.Order

	if err := json.Unmarshal(value, &order); err != nil {
		return fmt.Errorf("failed to unmarshal kafka message at offset %v: %w", offset, err)
	}

	if err := repo.Insert(context.Background(), &order); err != nil {
		return fmt.Errorf("failed to save order %s: %w", order.OrderUID, err)
	}

	return nil
}

func NewHandler() *Handler {
	return &Handler{}
}

func (r *Router) GetOrder(w http.ResponseWriter, req *http.Request) {
	log.Println("NEW REQUEST WITH ID")
	orderUid := chi.URLParam(req, "order_uid")
	if orderUid == "" {
		http.Error(w, "order_id is required", http.StatusBadRequest)
		return
	}

	order, err := r.service.GetOrder(orderUid)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(order)
	if err != nil {
		log.Println(err)
	}
	w.Write(resp)
}

func (r *Router) orderPage(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	tmplPath := filepath.Join("templates", "order.html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, "Ошибка рендеринга", http.StatusInternalServerError)
	}
}
