package storage

import (
	"sync"

	"github.com/google/uuid"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/model"
)

type GameStorage struct {
	games sync.Map
}

func NewGameStorage() *GameStorage {
	return &GameStorage{}
}

func (s *GameStorage) Store(game *model.GameEntity) {
	s.games.Store(game.ID, *game)
}

func (s *GameStorage) Load(id uuid.UUID) (*model.GameEntity, bool) {
	value, ok := s.games.Load(id)
	if !ok {
		return nil, false
	}
	game := value.(model.GameEntity)
	return &game, true
}
