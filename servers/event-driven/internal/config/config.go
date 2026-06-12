package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgreDSN  string
	KafkaBroker string
	RedisAddr   string
	GRPCPort    string
}

// Load data from .env
func Load() (*Config, error) {
	cfg := &Config{}

	if err := godotenv.Load(); err != nil {
		fmt.Println("no .env file, using system env vars")
	}

	dbHost := GetEnvFallback("DB_HOST", "localhost")
	dbPort := GetEnvFallback("DB_PORT", "5432")
	dbUser := GetEnvFallback("DB_USER", "postgres")
	dbName := GetEnvFallback("DB_NAME", "event_driven")

	dbPass, err := GetEnvWithout("DB_PASSWORD")
	if err != nil {
		return nil, fmt.Errorf("Load: %w", err)
	}

	cfg.PostgreDSN = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	cfg.RedisAddr = GetEnvFallback("REDIS_ADDR", "localhost:6379")
	cfg.KafkaBroker = GetEnvFallback("KAFKA_BROKERS", "localhost:9092")
	cfg.GRPCPort = GetEnvFallback("GRPC_PORT", "9090")

	return cfg, nil
}

// Loading data from the env with oportinity to make default value
func GetEnvFallback(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}

// Loading data from the env without default value
func GetEnvWithout(key string) (string, error) {
	if value := os.Getenv(key); value != "" {
		return value, nil
	}

	return "", fmt.Errorf("Can not load data from env with key: %s", key)
}
