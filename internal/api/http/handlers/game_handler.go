package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/dto"
	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/mapper"
	"github.com/jorikgrozniy/tic-tac-toe/internal/application"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/service"
)

type GameHandler struct {
	appService *application.GameServiceWithRepo
}

func NewGameHandler(appService *application.GameServiceWithRepo) *GameHandler {
	return &GameHandler{
		appService: appService,
	}
}

func (h *GameHandler) MakeMoveHandler(w http.ResponseWriter, r *http.Request, gameID string) {
	defer r.Body.Close()

	if gameID == "" {
		sendJSONerror(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	var req dto.MakeMoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendJSONerror(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	domainRequest := mapper.FromDTO(gameID, req)
	domainResponse, status, err := h.appService.ProcessMakeMove(domainRequest)
	if err != nil {
		sendJSONerror(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := mapper.ToDTO(domainResponse, status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendJSONerror(w, "Error creating response", http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	domainGame, err := h.appService.ProcessCreateGame()
	if err != nil {
		sendJSONerror(w, err.Error(), http.StatusInternalServerError)
	}

	response := mapper.ToDTO(domainGame, service.Playing)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendJSONerror(w, "Error creating response", http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) GetGameHandler(w http.ResponseWriter, r *http.Request, gameID string) {
	if gameID == "" {
		sendJSONerror(w, "Game ID is required", http.StatusBadRequest)
		return
	}

	domainResponse, status, err := h.appService.ProcessGetGame(gameID)
	if err != nil {
		sendJSONerror(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := mapper.ToDTO(domainResponse, status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		sendJSONerror(w, "Error creating response", http.StatusInternalServerError)
		return
	}
}

func sendJSONerror(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	errorResponse := map[string]string{"error": message}

	json.NewEncoder(w).Encode(errorResponse)
}
