package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Kafka           Kafka
	OrderRepository OrderRepository
}

type Service struct {
	Port string `env:"SERVICE_PORT"`
}

type OrderRepository struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
}

type Kafka struct {
	KafkaAddress  string `env:"KAFKA_ADDRESS"`
	Topic         string `env:"KAFKA_TOPIC"`
	ConsumerGroup string `env:"KAFKA_CONS_GROUP"`
}

func MustLoad() Config {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}
	return cfg
}
