package tests

import (
	"context"
	"event-driven/internal/config"
	"event-driven/internal/repository"
	"event-driven/internal/storage"
	"testing"
)

func TestVoteIntegration(t *testing.T) {
	cfg := &config.Config{PostgreDSN: "postgres://postgres:postgres@localhost:5432/event_driven_test?sslmode=disable"}
	db, err := storage.NewPostgresDB(cfg)
	if err != nil {
		t.Skip("Skipping integration test: no test DB")
	}
	defer db.Close()

	repo := repository.NewVoteRepository(db)
	err = repo.Cast(context.Background(), "poll-id", "option-id", "user1")
	if err != nil {
		t.Errorf("Cast failed: %v", err)
	}
}
