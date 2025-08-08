package api

import (
	"log"
	"net/http"
	"wb_tech_L0/internal/config"

	"github.com/go-chi/chi"
)

func NewRouter(cfg config.Service, service *service.Service, logger *log.Logger) *Router {
	innerRouter := chi.NewRouter()
	r := &Router{
		router: innerRouter,
		serv: &http.Server{
			Addr:    cfg.Port,
			Handler: innerRouter,
		},
		logger:  logger,
		service: service,
	}

	r.router.Get("/order/{order_uid}", r.getOrder)

	return r
}

type Router struct {
	serv    *http.Server
	router  *chi.Mux
	service *service.Service
	logger  *log.Logger
}
