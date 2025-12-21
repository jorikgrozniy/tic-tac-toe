package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/mapper"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/storage"
)

type GameRepositoryImpl struct {
	storage *storage.GameStorage
}

func NewGameRepositoryImpl(storage *storage.GameStorage) *GameRepositoryImpl {
	return &GameRepositoryImpl{
		storage: storage,
	}
}

func (r *GameRepositoryImpl) Save(game *model.CurrentGame) error {
	if game == nil {
		return errors.New("game is nil")
	}

	entity := mapper.ToDataSourceEntity(game)
	if entity.ID == uuid.Nil {
		entity.ID = uuid.New()
	}

	storedGame, exists := r.storage.Load(entity.ID)
	now := time.Now()
	if exists {
		entity.CreatedAt = storedGame.CreatedAt
	} else {
		entity.CreatedAt = now
	}
	entity.UpdatedAt = now

	r.storage.Store(entity)
	return nil
}

func (r *GameRepositoryImpl) Retrieve(id uuid.UUID) (*model.CurrentGame, error) {
	game, ok := r.storage.Load(id)
	if ok {
		return mapper.ToDomainEntity(game), nil
	} else {
		return nil, errors.New("game not found")
	}
}
