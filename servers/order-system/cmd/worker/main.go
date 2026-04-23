package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order-system/internal/config"
	"order-system/internal/kafka"
	"order-system/internal/models"
	"order-system/internal/storage"
)

func main() {
	op := "worker.main"
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	consumer := kafka.NewConsumer(cfg.KafkaBroker, cfg.OrdersTopic, cfg.ConsumerGroup)
	defer consumer.Close()

	logger.Info("Worker started")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := storage.NewDB(context.Background(), cfg.PostgresDSN)
	if err != nil {
		logger.Error("failed to create db", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}
	defer db.Close(context.Background())

	if err := db.Init(context.Background()); err != nil {
		logger.Error("failed to init db", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	for {
		if ctx.Err() != nil {
			logger.Info("Worker stopping...")
			break
		}

		msg, err := consumer.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				continue
			}

			logger.Error("failed to fetch message", slog.String("error", err.Error()), slog.String("op", op))
			time.Sleep(1 * time.Second)
			continue
		}

		var order models.OrderRequest
		if err = json.Unmarshal(msg.Value, &order); err != nil {
			logger.Error("failed to unmarshal order event", slog.String("error", err.Error()), slog.String("op", op))

			if err = consumer.CommitMessage(ctx, msg); err != nil {
				logger.Error("failed to commit invalid message", slog.String("error", err.Error()), slog.String("op", op))
				continue
			}

			continue
		}

		logger.Info("processing order", slog.String("id", order.ID), slog.String("op", op))

		if err := db.ProcessOrder(ctx, order); err != nil {
			logger.Error("failed to process order", slog.String("error", err.Error()), slog.String("op", op))
			time.Sleep(1 * time.Second)
			continue
		}

		if err = consumer.CommitMessage(ctx, msg); err != nil {
			logger.Error("failed to commit message", slog.String("error", err.Error()), slog.String("op", op))
		}

		logger.Info("order committed", slog.String("id", order.ID), slog.String("op", op))
	}
}
