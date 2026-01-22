package delete

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"url-shortener/internal/storage/postgres"

	"github.com/go-chi/chi/v5"
)

type Response struct {
	Status string
	Error  string
}

type DeleterUrl interface {
	DeleteUrl(ctx context.Context, alias string) error
}

func Delete(log *slog.Logger, deleterUrl DeleterUrl) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := chi.URLParam(r, "alias")

		if alias == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			msg := "bad link nothing to delete"
			log.Error("can't delete", "error", msg, "method", "Delete")
			json.NewEncoder(w).Encode(Response{
				Status: "Error",
				Error:  msg,
			})

			return
		}

		err := deleterUrl.DeleteUrl(r.Context(), alias)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")

			if errors.Is(err, postgres.ErrNotFound) {
				msg := "link not found"
				log.Warn(msg, "method", "Delete")
				json.NewEncoder(w).Encode(Response{
					Status: "Warn",
					Error:  msg,
				})

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("Delete: failed to delete url: %v", err)
			log.Error("failed to delete url", "error", err, "method", "Delete")
			json.NewEncoder(w).Encode(Response{
				Status: "Error",
				Error:  msg,
			})

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		log.Info("deleted url", "method", "Delete")
		json.NewEncoder(w).Encode(Response{
			Status: "OK",
			Error:  "",
		})
	}
}
