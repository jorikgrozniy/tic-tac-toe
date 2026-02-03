package handler

import (
	"encoding/json"
	"net/http"

	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/dto"
	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/util"
	"github.com/jorikgrozniy/tic-tac-toe/internal/auth"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/service"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/repository"
)

type AuthHandler struct {
	authService *auth.AuthService
}

func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	if err := h.authService.Register(&req); err != nil {
		switch err {
		case service.ErrUserAlreadyExists:
			util.SendError(w, err.Error(), http.StatusConflict)
		default:
			util.SendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	util.SendSuccess(w, "registration completed", http.StatusCreated)
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		util.SendError(w, "authorization header required", http.StatusUnauthorized)
		return
	}

	userID, err := h.authService.Authenticate(authHeader)
	if err != nil {
		switch err {
		case auth.ErrInvalidCredentials:
			util.SendError(w, err.Error(), http.StatusUnauthorized)
		case repository.ErrUserNotFound:
			util.SendError(w, err.Error(), http.StatusNotFound)
		default:
			util.SendError(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	util.SendAuthSuccess(w, userID.String())
}
