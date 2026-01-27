package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"order-system/internal/kafka"
	"order-system/internal/models"
	"order-system/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type OrderHandler struct {
	producer *kafka.Producer
	log      *slog.Logger
	db       *storage.DB
}

func NewOrderHandler(producer *kafka.Producer, log *slog.Logger, db *storage.DB) *OrderHandler {
	return &OrderHandler{
		producer: producer,
		log:      log,
		db:       db,
	}
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var JSONdata models.OrderRequest
	op := "handleCreateOrder"
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewDecoder(r.Body).Decode(&JSONdata); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := "json decode failed"
		h.log.Error("json decode failed", slog.String("error", err.Error()), slog.String("method", r.Method), slog.String("op", op))
		json.NewEncoder(w).Encode(models.OrderResponse{
			Status: "Error",
			Msg:    msg,
		})

		return
	}

	if err := JSONdata.Validate(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Sprintf("Validate data error: %v", err)
		h.log.Error("bad data", slog.String("error", msg), slog.String("method", r.Method), slog.String("op", op))
		json.NewEncoder(w).Encode(models.OrderResponse{
			Status: "Error",
			Msg:    msg,
		})

		return
	}

	bytes, err := json.Marshal(JSONdata)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("error with marshaling json: %v", err)
		h.log.Error("json marshal", slog.String("error", msg), slog.String("method", r.Method), slog.String("op", op))
		json.NewEncoder(w).Encode(models.OrderResponse{
			Status: "Error",
			Msg:    msg,
		})

		return
	}

	if err = h.producer.SendMessage(r.Context(), bytes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("Error with send message: %v", err)
		h.log.Error("Send message", slog.String("error", msg), slog.String("method", r.Method), slog.String("op", op))
		json.NewEncoder(w).Encode(models.OrderResponse{
			Status: "Error",
			Msg:    msg,
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	h.log.Info("message was send", slog.String("method", r.Method), slog.String("op", op))
	json.NewEncoder(w).Encode(models.OrderResponse{
		Status: "Accepted",
		Msg:    "",
	})
}

func (h *OrderHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	op := "GetOrderByID"
	id := chi.URLParam(r, "id")

	model, err := h.db.GetOrder(r.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			msg := fmt.Sprintf("Error with bad data: %v", err)
			h.log.Warn("Bad data", slog.String("warn", err.Error()), slog.String("method", r.Method), "op", op)
			json.NewEncoder(w).Encode(models.OrderResponse{
				Status: "Warn",
				Msg:    msg,
			})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("Error with get data: %v", err)
		h.log.Error("Error with get data", slog.String("error", msg), slog.String("method", r.Method), slog.String("op", op))
		json.NewEncoder(w).Encode(models.OrderResponse{
			Status: "Error",
			Msg:    msg,
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	h.log.Info("order was got", slog.String("method", r.Method), slog.String("op", op))
	json.NewEncoder(w).Encode(models.OrderRequest{
		ID:    model.ID,
		Item:  model.Item,
		Price: model.Price,
	})

}

func (h *OrderHandler) DeleteOrderByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")
	op := "DeleteOrderByID"

	if err := h.db.DeleteOrder(r.Context(), id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			w.WriteHeader(http.StatusNotFound)
			msg := fmt.Sprintf("Order was not found: %v", err)
			h.log.Warn("order by ID was not found", slog.String("warn", err.Error()), slog.String("method", r.Method), slog.String("op", op))
			json.NewEncoder(w).Encode(models.OrderResponse{
				Status: "Warn",
				Msg:    msg,
			})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Sprintf("Error with delete data: %v", err)
		h.log.Error("failed to delete data", slog.String("warn", err.Error()), slog.String("method", r.Method), slog.String("op", op))
		json.NewEncoder(w).Encode(models.OrderResponse{
			Status: "Error",
			Msg:    msg,
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	h.log.Info("order was deleted", slog.String("op", op))
	json.NewEncoder(w).Encode(models.OrderResponse{
		Status: "Deleted",
		Msg:    "",
	})
}
