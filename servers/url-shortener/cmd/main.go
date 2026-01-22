package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/storage/postgres"
	"url-shortener/internal/storage/redis"
	storageservice "url-shortener/internal/storage/storage_service"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.LoadConfig()

	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	log.Info("starting url-shortener", slog.String("env", "local"))

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	pgClient, err := postgres.NewDB(dsn)
	if err != nil {
		log.Error("failed to connect to postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pgClient.Close()

	log.Info("connected to storage")

	redisClient, err := redis.NewClient(cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		log.Error("failed to connect to redis", slog.String("error", err.Error()))
		os.Exit(1)
	}

	defer redisClient.Close()

	log.Info("connected to redis")

	storageService := storageservice.NewService(log, pgClient, redisClient)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Post("/url", save.New(log, storageService))
	router.Get("/{alias}", redirect.New(log, storageService))
	router.Delete("/delete/{alias}", delete.Delete(log, storageService))

	srv := &http.Server{
		Addr:    cfg.AppPort,
		Handler: router,
	}

	go func() {
		log.Info("starting server", slog.String("address", cfg.AppPort))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("failed to start server", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Info("Server forced to shutdown:", "error", err)
	}

	log.Info("Server exiting")
}
