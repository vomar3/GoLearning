package handler

import (
	apperrors "auth-service/errors"
	jwtlib "auth-service/internal/jwt"
	"auth-service/internal/models"
	"auth-service/internal/storage"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type AuthToken struct {
	Token string `json:"token"`
}

type AuthHandler struct {
	db     *storage.DB
	logger *slog.Logger
}

func NewAuthHandler(log *slog.Logger, db *storage.DB) *AuthHandler {
	return &AuthHandler{
		logger: log,
		db:     db,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	op := "handlers.Register"
	w.Header().Set("Content-Type", "application/json")
	var NewUser models.User

	if err := json.NewDecoder(r.Body).Decode(&NewUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Errorf("Failed decode json: %w", err)
		h.logger.Error("Failed to decode json", slog.String("error", err.Error()), slog.String("op", op))
		json.NewEncoder(w).Encode(models.UserBadResponse{
			Error: "error",
			Msg:   msg.Error(),
		})

		return
	}

	if NewUser.Validate() == false {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Errorf("Failed to validate data, check your data and try again")
		h.logger.Error("Failed to validate data", slog.String("error", msg.Error()), slog.String("op", op))
		json.NewEncoder(w).Encode(models.UserBadResponse{
			Error: "error",
			Msg:   msg.Error(),
		})

		return
	}

	id, err := h.db.CreateUser(r.Context(), NewUser.Email, NewUser.Password)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			msg := fmt.Errorf("The user already exists: %w", err)
			h.logger.Warn("The user already exists", slog.String("warn", err.Error()), slog.String("op", op))
			json.NewEncoder(w).Encode(models.UserBadResponse{
				Error: "warn",
				Msg:   msg.Error(),
			})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Errorf("Failed to create user: %w", err)
		h.logger.Error("Failed to create user", slog.String("error", err.Error()), slog.String("op", op))
		json.NewEncoder(w).Encode(models.UserBadResponse{
			Error: "error",
			Msg:   msg.Error(),
		})

		return
	}

	w.WriteHeader(http.StatusCreated)
	h.logger.Info("The user was successfully registered", slog.String("op", op))
	json.NewEncoder(w).Encode(models.UserGoodResponse{
		UserID: id,
	})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	op := "auth.Login"
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Errorf("Failed decode json: %w", err)
		h.logger.Error("Failed to decode json", slog.String("error", err.Error()), slog.String("op", op))
		json.NewEncoder(w).Encode(models.UserBadResponse{
			Error: "error",
			Msg:   msg.Error(),
		})

		return
	}

	userData, err := h.db.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		if errors.Is(err, apperrors.ErrUserNotFound) {
			w.WriteHeader(http.StatusBadRequest)
			msg := fmt.Errorf("Invalid email or password: %w", err)
			h.logger.Warn("User was not found in db", slog.String("warn", err.Error()), slog.String("op", op))
			json.NewEncoder(w).Encode(models.UserBadResponse{
				Error: "warn",
				Msg:   msg.Error(),
			})

			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Errorf("Invalid email or password: %w", err)
		h.logger.Error("Failed to found user", slog.String("error", err.Error()), slog.String("op", op))
		json.NewEncoder(w).Encode(models.UserBadResponse{
			Error: "error",
			Msg:   msg.Error(),
		})

		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userData.PasswordHash), []byte(user.Password)); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		msg := fmt.Errorf("bad data: %w", err)
		h.logger.Warn("bad data", slog.String("warn", err.Error()), slog.String("op", op))
		json.NewEncoder(w).Encode(models.UserBadResponse{
			Error: "warn",
			Msg:   msg.Error(),
		})

		return
	}

	token, err := jwtlib.GenerateToken(userData.UserID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		msg := fmt.Errorf("Failed to create token: %w", err)
		h.logger.Error("Failed to create token", slog.String("error", err.Error()), slog.String("op", op))
		json.NewEncoder(w).Encode(models.UserBadResponse{
			Error: "error",
			Msg:   msg.Error(),
		})

		return
	}

	w.WriteHeader(http.StatusOK)
	h.logger.Info("token was generated")
	json.NewEncoder(w).Encode(AuthToken{
		Token: token,
	})
}
