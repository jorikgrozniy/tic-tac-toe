package repository

import (
	"tic-tac-toe/internal/domain/model"

	"github.com/google/uuid"
)

type GameRepository interface {
	Save(game *model.CurrentGame) error
	Retrieve(id uuid.UUID) (*model.CurrentGame, error)
	FindAvailableGames() ([]model.GameInfo, error)
}
