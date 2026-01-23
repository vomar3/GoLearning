package redirect

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Url    string `json:"long_url"`
}

type URLGetter interface {
	GetUrl(ctx context.Context, alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.redirect.New"
		alias := chi.URLParam(r, "alias")

		if alias == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			msg := "bad alias"
			log.Error("0 length alias", slog.String("error", msg), slog.String("method", r.Method), slog.String("op", op))
			json.NewEncoder(w).Encode(Response{
				Status: "Error",
				Error:  msg,
				Url:    "",
			})

			return
		}

		url, err := urlGetter.GetUrl(r.Context(), alias)

		if err != nil {
			w.Header().Set("Content-Type", "application/json")

			if errors.Is(err, pgx.ErrNoRows) {
				w.WriteHeader(http.StatusNotFound)
				msg := fmt.Sprintf("Url wasn't found: %v", err)
				log.Error("Url wasn't found", slog.String("error", err.Error()), slog.String("method", r.Method), slog.String("op", op))
				json.NewEncoder(w).Encode(Response{
					Status: "Error",
					Error:  msg,
					Url:    "",
				})

				return
			}

			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("failed to get URL: %v", err)
			log.Error("failed to get URL", slog.String("error", err.Error()), slog.String("method", r.Method), slog.String("op", op))
			json.NewEncoder(w).Encode(Response{
				Status: "Error",
				Error:  msg,
				Url:    "",
			})

			return
		}

		log.Info("URL for %s alias was found: %s", alias, url)
		http.Redirect(w, r, url, http.StatusFound)
	}
}
