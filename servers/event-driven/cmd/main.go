package main

import (
	"event-driven/internal/config"
	"event-driven/internal/storage"
	"log/slog"
	"os"
)

func main() {
	op := "main"
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.Load()
	if err != nil {
		logger.Error("failed to load config", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	db, err := storage.NewPostgresDB(cfg)
	if err != nil {
		logger.Error("failed to connect to db", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("Connect to db was successful", slog.String("op", op))
}
