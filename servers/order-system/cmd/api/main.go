package main

import (
	"context"
	"log/slog"
	"net/http"
	"order-system/cmd/api/handlers"
	"order-system/internal/kafka"
	"order-system/internal/storage"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type OrderRequest struct {
	ID    string `json:"id"`
	Item  string `json:"item"`
	Price int    `json:"price"`
}

type OrderResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg,omitempty"`
}

func main() {
	op := "api.main"
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	producer := kafka.NewProducer("localhost:9092", "orders")
	defer producer.Close()

	db, err := storage.NewDB(context.Background(), "postgres://user:password@localhost:5444/orders")
	if err != nil {
		logger.Error("failed to connect to db", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	defer db.Close(context.Background())

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	srv := &http.Server{
		Addr:    ":8080", // Надо конфиг позже написать
		Handler: router,
	}

	myHandler := handlers.NewOrderHandler(producer, logger, db)

	router.Post("/order", myHandler.CreateOrder)
	router.Get("/orders/{id}", myHandler.GetOrderByID)
	router.Delete("/orders/{id}", myHandler.DeleteOrderByID)

	logger.Info("starting server", slog.String("port", "8080"))
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("main: error with statring server", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Info("Server forced to shutdown:", "error", err)
	}

	logger.Info("Server was closed")
}
