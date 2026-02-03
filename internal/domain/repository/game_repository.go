package repository

import (
	"github.com/google/uuid"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
)

type GameRepository interface {
	Save(game *model.CurrentGame) error
	Retrieve(id uuid.UUID) (*model.CurrentGame, error)
	FindAvailableGames() ([]model.GameInfo, error)
}
