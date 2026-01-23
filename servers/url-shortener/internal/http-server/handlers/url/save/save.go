package save

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"url-shortener/lib/random"
)

type Request struct {
	Url   string `json:"long_url"`
	Alias string `json:"short_code"`
}

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error"`
	Alias  string `json:"short_code"`
}

type URLSaver interface {
	SaveUrl(ctx context.Context, url string, alias string) error
}

func (r *Request) CheckAlias() string {
	if r.Alias == "" {
		r.Alias = random.Rand(6)
	}

	return r.Alias
}

func (r *Request) Validate() bool {
	if r.Url == "" {
		return false
	}

	countDots := 0
	for _, val := range r.Url {
		if val == '.' {
			countDots++
		}
	}

	if countDots < 1 {
		return false
	}

	return true
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		var data Request

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("Can't unmarshal data from json: %v", err)
			log.Error("json decode failed", slog.String("error", msg), slog.String("method", r.Method), slog.String("op", op))
			json.NewEncoder(w).Encode(Response{
				Status: "Error",
				Error:  msg,
				Alias:  "",
			})

			return
		}

		if data.Validate() != true {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			msg := "bad validate data"
			log.Error("validate failed", slog.String("error", msg), slog.String("method", r.Method), slog.String("op", op))
			json.NewEncoder(w).Encode(Response{
				Status: "Error",
				Error:  msg,
				Alias:  "",
			})

			return
		}

		alias := data.CheckAlias()

		if err := urlSaver.SaveUrl(r.Context(), data.Url, alias); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("error with save data: %v", err)
			log.Error("save url failed", slog.String("error", msg), slog.String("method", r.Method), slog.String("op", op))
			json.NewEncoder(w).Encode(Response{
				Status: "Error",
				Error:  msg,
				Alias:  "",
			})

			return
		}

		if _, err := json.Marshal(data); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("error with marshal data: %v", err)
			slog.Error("json marshal failed", slog.String("error", msg), slog.String("method", r.Method), slog.String("op", op))
			json.NewEncoder(w).Encode(Response{
				Status: "Error",
				Error:  msg,
				Alias:  "",
			})

			return
		}

		log.Info("url saved", slog.String("alias", alias), slog.String("op", op))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(Response{
			Status: "OK",
			Error:  "",
			Alias:  alias,
		})
	}
}
