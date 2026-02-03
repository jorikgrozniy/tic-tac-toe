package service

import (
	"github.com/google/uuid"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
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
