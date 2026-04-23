package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"order-system/internal/models"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type fakeProducer struct {
	key     string
	message []byte
	err     error
}

func (p *fakeProducer) SendMessage(_ context.Context, key string, message []byte) error {
	p.key = key
	p.message = message
	return p.err
}

type fakeStore struct {
	order models.Order
	err   error
}

func (s *fakeStore) GetOrder(_ context.Context, _ string) (models.Order, error) {
	return s.order, s.err
}

func (s *fakeStore) DeleteOrder(_ context.Context, _ string) error {
	return s.err
}

func TestCreateOrderSendsKafkaMessage(t *testing.T) {
	producer := &fakeProducer{}
	handler := NewOrderHandler(producer, testLogger(), &fakeStore{})

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBufferString(`{"id":"order-1","item":"book","price":1200}`))
	rec := httptest.NewRecorder()

	handler.CreateOrder(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}
	if producer.key != "order-1" {
		t.Fatalf("expected Kafka key order-1, got %q", producer.key)
	}

	var sent models.OrderRequest
	if err := json.Unmarshal(producer.message, &sent); err != nil {
		t.Fatalf("failed to unmarshal sent message: %v", err)
	}
	if sent.ID != "order-1" || sent.Item != "book" || sent.Price != 1200 {
		t.Fatalf("unexpected sent order: %+v", sent)
	}
}

func TestCreateOrderRejectsInvalidPayload(t *testing.T) {
	producer := &fakeProducer{}
	handler := NewOrderHandler(producer, testLogger(), &fakeStore{})

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBufferString(`{"id":"","item":"book","price":1200}`))
	rec := httptest.NewRecorder()

	handler.CreateOrder(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", rec.Code)
	}
	if producer.message != nil {
		t.Fatal("expected producer not to be called")
	}
}

func TestGetOrderByIDReturnsOrder(t *testing.T) {
	now := time.Date(2026, 4, 23, 12, 0, 0, 0, time.UTC)
	store := &fakeStore{
		order: models.Order{
			ID:        "order-1",
			Item:      "book",
			Price:     1200,
			Status:    models.OrderStatusProcessed,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
	handler := NewOrderHandler(&fakeProducer{}, testLogger(), store)

	req := httptest.NewRequest(http.MethodGet, "/orders/order-1", nil)
	rec := httptest.NewRecorder()
	handler.GetOrderByID(rec, withURLParam(req, "id", "order-1"))

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var got models.Order
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Status != models.OrderStatusProcessed {
		t.Fatalf("expected processed status, got %q", got.Status)
	}
}

func TestDeleteOrderByIDReturnsNotFound(t *testing.T) {
	handler := NewOrderHandler(&fakeProducer{}, testLogger(), &fakeStore{err: pgx.ErrNoRows})

	req := httptest.NewRequest(http.MethodDelete, "/orders/order-404", nil)
	rec := httptest.NewRecorder()
	handler.DeleteOrderByID(rec, withURLParam(req, "id", "order-404"))

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}

func TestCreateOrderReturnsServerErrorWhenProducerFails(t *testing.T) {
	handler := NewOrderHandler(&fakeProducer{err: errors.New("kafka down")}, testLogger(), &fakeStore{})

	req := httptest.NewRequest(http.MethodPost, "/order", bytes.NewBufferString(`{"id":"order-1","item":"book","price":1200}`))
	rec := httptest.NewRecorder()

	handler.CreateOrder(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", rec.Code)
	}
}

func withURLParam(req *http.Request, key, value string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)

	ctx := context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx)
	return req.WithContext(ctx)
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}
