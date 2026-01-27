package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"order-system/internal/kafka"
	"order-system/internal/models"
	"order-system/internal/storage"
)

func main() {
	op := "worker.main"
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	consumer := kafka.NewConsumer("localhost:9092", "orders", "order-workers")
	defer consumer.Close()

	logger.Info("Worker started")

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	db, err := storage.NewDB(context.Background(),
		"postgres://user:password@localhost:5444/orders")

	if err != nil {
		logger.Error("failed to create bd", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	defer db.Close(context.Background())

	if err := db.Init(context.Background()); err != nil {
		logger.Error("failed to create bd", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	for {
		if ctx.Err() != nil {
			logger.Info("Worker stopping...")
			break
		}

		// Блокирующаяся операция
		msg, err := consumer.FetchMessage(ctx)
		if err != nil {
			if ctx.Err() != nil {
				continue
			}

			logger.Error("failed to Fetch message", slog.String("error", err.Error()), slog.String("op", op))
			time.Sleep(1 * time.Second)
			continue
		}

		var order models.OrderRequest
		if err = json.Unmarshal(msg.Value, &order); err != nil {
			logger.Error("failed to Unmarshal", slog.String("error", err.Error()), slog.String("op", op))

			// Заливаем, чтобы "битое" сообщение не попадалось нам из разу в раз
			if err = consumer.CommitMessage(ctx, msg); err != nil {
				logger.Error("failed to Commit message", slog.String("error", err.Error()), slog.String("op", op))
				continue
			}

			continue
		}

		logger.Info("Processing order...", slog.String("id", order.ID), slog.String("op", op))

		if err := db.SaveOrder(ctx, order); err != nil {
			logger.Error("failed to save order", slog.String("error", err.Error()), slog.String("op", op))

			time.Sleep(1 * time.Second)
			continue
		}

		if err = consumer.CommitMessage(ctx, msg); err != nil {
			logger.Error("failed to Commit message", slog.String("error", err.Error()), slog.String("op", op))
		}

		logger.Info("Order committed", slog.String("id", order.ID), slog.String("op", op))
	}
}
