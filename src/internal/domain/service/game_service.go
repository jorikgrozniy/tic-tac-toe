package service

import (
	"tic-tac-toe/internal/domain/model"

	"github.com/google/uuid"
)

type GameService interface {
	ComputeNextMove(game *model.CurrentGame, player model.GamePlayer) (int, int, error)
	ValidateGameBoard(game *model.CurrentGame, player model.GamePlayer) error
	CheckGameCompletion(game *model.CurrentGame) model.GameStatus

	MakeMove(reqGame *model.CurrentGame, userID uuid.UUID) error
	CreateGame(userID uuid.UUID, t string) (*model.CurrentGame, error)
	GetGame(gameID string, userID uuid.UUID) (*model.CurrentGame, error)
	GetAvailableGames() ([]model.GameInfo, int, error)
	JoinGame(gameID string, userID uuid.UUID) error
}
