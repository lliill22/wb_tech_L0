package handlers

import (
	"context"
	"log"
	"net/http"
	"wb_tech_L0/internal/config"
	"wb_tech_L0/internal/service"

	"github.com/go-chi/chi"
)

type Router struct {
	server  *http.Server
	chi     *chi.Mux
	service *service.Service
}

func (r *Router) Start() {
	if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}

func (r *Router) Close(ctx context.Context) error {
	return r.server.Shutdown(ctx)
}

func NewRouter(cfg config.Service, serv *service.Service) *Router {
	innerRouter := chi.NewRouter()
	r := &Router{
		server: &http.Server{
			Addr:    cfg.Port,
			Handler: innerRouter,
		},
		chi:     innerRouter,
		service: serv,
	}

	r.chi.Get("/order/{order_uid}", r.GetOrder)
	r.chi.Get("/order/", r.orderPage)
	return r
}
