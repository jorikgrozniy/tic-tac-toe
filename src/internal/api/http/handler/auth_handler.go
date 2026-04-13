package handler

import (
	"encoding/json"
	"net/http"

	"tic-tac-toe/internal/api/http/contextkey"
	"tic-tac-toe/internal/api/http/dto"
	"tic-tac-toe/internal/api/http/mapper"
	"tic-tac-toe/internal/api/http/util"
	"tic-tac-toe/internal/auth"
	"tic-tac-toe/internal/domain/service"
	"tic-tac-toe/internal/infrastructure/datasource/repository"
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
	var req auth.JwtRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	jwtResponse, err := h.authService.Authenticate(&req)
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

	if err := json.NewEncoder(w).Encode(*jwtResponse); err != nil {
		util.SendError(w, "error creating response", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) UpdateAccessTokenHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	jwtResponse, err := h.authService.UpdateAccessToken(authHeader)
	if err != nil {
		util.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(*jwtResponse); err != nil {
		util.SendError(w, "error creating response", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) UpdateRefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	jwtResponse, err := h.authService.UpdateRefreshToken(authHeader)
	if err != nil {
		util.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(*jwtResponse); err != nil {
		util.SendError(w, "error creating response", http.StatusInternalServerError)
		return
	}
}

func (h *AuthHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := contextkey.GetUserID(r.Context())
	if !ok {
		util.SendError(w, "Authentication required", http.StatusUnauthorized)
		return
	}

	userInfo, err := h.authService.GetMyInfo(userID.String())
	if err != nil {
		util.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := mapper.ToDTOUserInfo(userInfo)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		util.SendError(w, "error creating response", http.StatusInternalServerError)
		return
	}
}
