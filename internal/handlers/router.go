package handlers

import (
	"context"
	"net/http"
	"wb_tech_L0/internal/config"

	"github.com/go-chi/chi"
)

type Router struct {
	server *http.Server
	chi    *chi.Mux
}

func (r *Router) Start() {
	r.server.ListenAndServe()
}

func (r *Router) Close() {
	r.server.Shutdown(context.Background())
}

func NewRouter(cfg config.Service) *Router {
	innerRouter := chi.NewRouter()
	r := &Router{
		server: &http.Server{
			Addr:    cfg.Port,
			Handler: innerRouter,
		},
		chi: innerRouter,
	}

	r.chi.Get("/{order_id}", r.GetOrder)
	return r
}
