package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	HTTPPort        string
	PostgresDSN     string
	KafkaBroker     string
	OrdersTopic     string
	ConsumerGroup   string
	ShutdownTimeout time.Duration
}

func Load() Config {
	return Config{
		HTTPPort:        getEnv("HTTP_PORT", "8080"),
		PostgresDSN:     getEnv("POSTGRES_DSN", "postgres://user:password@localhost:5444/orders"),
		KafkaBroker:     getEnv("KAFKA_BROKER", "localhost:9092"),
		OrdersTopic:     getEnv("KAFKA_ORDERS_TOPIC", "orders"),
		ConsumerGroup:   getEnv("KAFKA_CONSUMER_GROUP", "order-workers"),
		ShutdownTimeout: time.Duration(getEnvInt("SHUTDOWN_TIMEOUT_SECONDS", 5)) * time.Second,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

func getEnvInt(key string, fallback int) int {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}

	return parsed
}
