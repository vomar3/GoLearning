package config

import (
	"testing"
	"time"
)

func TestLoadUsesDefaults(t *testing.T) {
	t.Setenv("HTTP_PORT", "")
	t.Setenv("POSTGRES_DSN", "")
	t.Setenv("KAFKA_BROKER", "")
	t.Setenv("KAFKA_ORDERS_TOPIC", "")
	t.Setenv("KAFKA_CONSUMER_GROUP", "")
	t.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "")

	cfg := Load()

	if cfg.HTTPPort != "8080" {
		t.Fatalf("expected default HTTP port 8080, got %q", cfg.HTTPPort)
	}
	if cfg.KafkaBroker != "localhost:9092" {
		t.Fatalf("expected default Kafka broker, got %q", cfg.KafkaBroker)
	}
	if cfg.ShutdownTimeout != 5*time.Second {
		t.Fatalf("expected default shutdown timeout 5s, got %s", cfg.ShutdownTimeout)
	}
}

func TestLoadUsesEnvironment(t *testing.T) {
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("POSTGRES_DSN", "postgres://test:test@localhost:5432/test")
	t.Setenv("KAFKA_BROKER", "kafka:29092")
	t.Setenv("KAFKA_ORDERS_TOPIC", "orders.test")
	t.Setenv("KAFKA_CONSUMER_GROUP", "workers-test")
	t.Setenv("SHUTDOWN_TIMEOUT_SECONDS", "12")

	cfg := Load()

	if cfg.HTTPPort != "9090" {
		t.Fatalf("expected env HTTP port, got %q", cfg.HTTPPort)
	}
	if cfg.PostgresDSN != "postgres://test:test@localhost:5432/test" {
		t.Fatalf("expected env Postgres DSN, got %q", cfg.PostgresDSN)
	}
	if cfg.KafkaBroker != "kafka:29092" {
		t.Fatalf("expected env Kafka broker, got %q", cfg.KafkaBroker)
	}
	if cfg.OrdersTopic != "orders.test" {
		t.Fatalf("expected env orders topic, got %q", cfg.OrdersTopic)
	}
	if cfg.ConsumerGroup != "workers-test" {
		t.Fatalf("expected env consumer group, got %q", cfg.ConsumerGroup)
	}
	if cfg.ShutdownTimeout != 12*time.Second {
		t.Fatalf("expected env shutdown timeout 12s, got %s", cfg.ShutdownTimeout)
	}
}
