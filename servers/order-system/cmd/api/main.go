package main

import (
	"context"
	"log/slog"
	"net/http"
	"order-system/cmd/api/handlers"
	"order-system/internal/config"
	"order-system/internal/kafka"
	"order-system/internal/storage"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	op := "api.main"
	cfg := config.Load()
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	producer := kafka.NewProducer(cfg.KafkaBroker, cfg.OrdersTopic)
	defer producer.Close()

	db, err := storage.NewDB(context.Background(), cfg.PostgresDSN)
	if err != nil {
		logger.Error("failed to connect to db", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}
	defer db.Close(context.Background())

	if err := db.Init(context.Background()); err != nil {
		logger.Error("failed to init db", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	srv := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	myHandler := handlers.NewOrderHandler(producer, logger, db)

	router.Post("/order", myHandler.CreateOrder)
	router.Get("/orders/{id}", myHandler.GetOrderByID)
	router.Delete("/orders/{id}", myHandler.DeleteOrderByID)

	logger.Info("starting server", slog.String("port", cfg.HTTPPort))
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("main: error with starting server", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Info("Server forced to shutdown:", "error", err)
	}

	logger.Info("Server was closed")
}
