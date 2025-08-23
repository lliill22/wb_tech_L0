package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"wb_tech_L0/internal/config"
	"wb_tech_L0/internal/handlers"
	kafka "wb_tech_L0/internal/kafka/consumer"
	"wb_tech_L0/internal/service"
	"wb_tech_L0/internal/storage"
)

func main() {
	conf := config.MustLoad()
	h := handlers.NewHandler()
	c, err := kafka.NewConsumer(h, conf)
	if err != nil {
		log.Fatal(err)
	}

	cache := storage.NewCache()
	db, err := storage.NewOrderRepository(conf.OrderRepository)
	if err != nil {
		log.Fatal(err)
	}
	service := service.NewService(&service.ServiceConfig{
		Cache:         cache,
		KafkaConsumer: c,
		OrderRepo:     db,
	})

	service.StartConsumer()
	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Fatal(service.StopConsumer().Error())

}
