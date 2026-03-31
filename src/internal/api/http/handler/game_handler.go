package handler

import (
	"encoding/json"
	"net/http"

	contextkeys "tic-tac-toe/internal/api/http/contextkey"
	"tic-tac-toe/internal/api/http/dto"
	"tic-tac-toe/internal/api/http/mapper"
	"tic-tac-toe/internal/api/http/util"
	"tic-tac-toe/internal/domain/service"
)

type GameHandler struct {
	gameService service.GameService
}

func NewGameHandler(gameService service.GameService) *GameHandler {
	return &GameHandler{
		gameService: gameService,
	}
}

func (h *GameHandler) MakeMoveHandler(w http.ResponseWriter, r *http.Request, gameID string) {
	if gameID == "" {
		util.SendError(w, "game ID is required", http.StatusBadRequest)
		return
	}

	userID, _ := contextkeys.GetUserID(r.Context())

	var req dto.MakeMoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	domainGame := mapper.FromDTOGame(gameID, req)
	err := h.gameService.MakeMove(domainGame, userID)
	if err != nil {
		util.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := mapper.ToDTOGame(domainGame)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		util.SendError(w, "error creating response", http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	userID, _ := contextkeys.GetUserID(r.Context())

	var req dto.CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		util.SendError(w, "invalid JSON format", http.StatusBadRequest)
		return
	}

	domainGame, err := h.gameService.CreateGame(userID, req.Type)
	if err != nil {
		util.SendError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response := mapper.ToDTOGame(domainGame)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		util.SendError(w, "error creating response", http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) GetGameHandler(w http.ResponseWriter, r *http.Request, gameID string) {
	if gameID == "" {
		util.SendError(w, "game ID is required", http.StatusBadRequest)
		return
	}

	userID, _ := contextkeys.GetUserID(r.Context())

	domainGame, err := h.gameService.GetGame(gameID, userID)
	if err != nil {
		util.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := mapper.ToDTOGame(domainGame)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		util.SendError(w, "error creating response", http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) GetAvailableGamesHandler(w http.ResponseWriter, r *http.Request) {
	domainGames, total, err := h.gameService.GetAvailableGames()
	if err != nil {
		util.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := mapper.ToDTOAvailableGames(domainGames, total)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		util.SendError(w, "error creating response", http.StatusInternalServerError)
		return
	}
}

func (h *GameHandler) JoinGameHandler(w http.ResponseWriter, r *http.Request, gameID string) {
	if gameID == "" {
		util.SendError(w, "game ID is required", http.StatusBadRequest)
		return
	}

	userID, _ := contextkeys.GetUserID(r.Context())

	err := h.gameService.JoinGame(gameID, userID)
	if err != nil {
		util.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	util.SendSuccess(w, "join successfull", http.StatusOK)
}
