package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Service  Service
	Postgres Postgres
	Kafka    Kafka
	Logger   log.Logger
}

type Service struct {
	Port string `env:"SERVICE_PORT"`
	Name string `env:"SERVICE_NAME"`
}

type Postgres struct {
	User     string `env:"POSTGRES_USER"`
	Password string `env:"POSTGRES_PASSWORD"`
	Database string `env:"POSTGRES_DB"`
	Host     string `env:"POSTGRES_HOST"`
	Port     string `env:"POSTGRES_PORT"`
}

type Kafka struct {
	Host      string `env:"KAFKA_HOST"`
	Port      string `env:"KAFKA_PORT"`
	UserTopic string `env:"USER_POST_CREATED"`
}

func MustLoad() *Config {
	cfg := &Config{}
	err := cleanenv.ReadEnv(cfg)
	if err != nil {
		log.Fatalf("Can not read env variables: %s", err)
	}
	return cfg
}
