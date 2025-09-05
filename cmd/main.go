package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"wb_tech_L0/internal/config"
	"wb_tech_L0/internal/handlers"
	kafka "wb_tech_L0/internal/kafka/consumer"
	"wb_tech_L0/internal/service"
	"wb_tech_L0/internal/storage"
)

func main() {
	cfg := config.MustLoad()
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	h := handlers.NewHandler()

	repo, err := storage.NewOrderRepository(cfg.OrderRepository)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	cache, err := storage.NewCache(context.Background(), repo)
	if err != nil {
		log.Fatal(err)
	}

	c, err := kafka.NewConsumer(h, cfg)
	if err != nil {
		log.Fatal(err)
	}

	serviceCfg := &service.ServiceConfig{
		Cache:         cache,
		KafkaConsumer: c,
		OrderRepo:     repo,
	}

	srv := service.NewService(serviceCfg)
	srv.StartConsumer()

	r := handlers.NewRouter(cfg.Service, srv)
	go r.Start()

	sig := <-sigCh
	log.Printf("got signal: %s, shutting down...", sig)

	srv.StopConsumer()

	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := r.Close(shutdownCtx); err != nil {
		log.Printf("http shutdown error: %v", err)
	}
	log.Println("shutdown complete")
}
