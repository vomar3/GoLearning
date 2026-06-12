package main

import (
	handler "auth-service/internal/handlers"
	myMiddleware "auth-service/internal/middleware"
	"auth-service/internal/storage"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	op := "main"
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	db, err := storage.NewDB(context.Background(), "postgres://user:password@localhost:5444/auth?sslmode=disable")
	if err != nil {
		logger.Error("Failed to create db", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	defer db.Close()

	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	myHandler := handler.NewAuthHandler(logger, db)

	router.Post("/register", myHandler.Register)
	router.Post("/login", myHandler.Login)

	router.Group(func(r chi.Router) {
		r.Use(myMiddleware.Auth)

		r.Get("/profile", func(w http.ResponseWriter, r *http.Request) {
			userID := r.Context().Value("userID").(int)

			w.Write([]byte(fmt.Sprintf("Hello, User %d! You are authorized.", userID)))
		})
	})

	logger.Info("Starting a server")
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Error with started the server", slog.String("error", err.Error()), slog.String("op", op))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Server is shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Error with shut down server")
	}

	logger.Info("Server was closed")
}
