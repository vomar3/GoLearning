package main

import (
	"context"
	"database/sql"
	"event-driven/internal/api/grpc"
	"event-driven/internal/config"
	"event-driven/internal/redis"
	"event-driven/internal/repository"
	"event-driven/internal/storage"
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	option "event-driven/proto/option"
	poll "event-driven/proto/poll"
	"event-driven/proto/vote"

	"github.com/joho/godotenv"
	grpcserver "google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	op := "main"
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}

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

	sqlDB, err := sql.Open("pgx", cfg.PostgreDSN)
	if err != nil {
		logger.Error("failed to open sql connection for migrations", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	defer sqlDB.Close()

	err = storage.RunMigrations(sqlDB, "./migrations")
	if err != nil {
		logger.Error("failed to migrate db", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	rds, err := redis.NewRedisClient(cfg, ctx)
	if err != nil {
		logger.Error("failed to connect to redis", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	logger.Info("connection to redis was successfully", slog.String("op", op))
	defer rds.Close()

	repo := repository.NewPollRepository(db)
	pollServer := grpc.NewPollServer(repo, rds)
	// create gRPC-server
	grpcServer := grpcserver.NewServer()
	// registration
	poll.RegisterPollServiceServer(grpcServer, pollServer)

	optionRepo := repository.NewOptionRepository(db)
	optionServer := grpc.NewOptionServer(optionRepo)
	option.RegisterOptionServiceServer(grpcServer, optionServer)

	voteRepo := repository.NewVoteRepository(db)
	voteServer := grpc.NewVoteServer(voteRepo, optionRepo, rds)
	vote.RegisterVoteServiceServer(grpcServer, voteServer)

	// for grpcurl
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		logger.Error("failed to listen", slog.String("error", err.Error()), slog.String("op", op))
		os.Exit(1)
	}

	go func() {
		logger.Info("gRPC server is running")
		if err := grpcServer.Serve(listener); err != nil {
			logger.Error("failed to serve", slog.String("error", err.Error()), slog.String("op", op))
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down gracefully...")
	grpcServer.GracefulStop()
	logger.Info("Server stopped")
}
