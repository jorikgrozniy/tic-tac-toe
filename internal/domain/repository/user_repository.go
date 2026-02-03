package repository

import (
	"github.com/google/uuid"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
)

type UserRepository interface {
	Save(game *model.User) error
	FindByID(id uuid.UUID) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindUserInfo(id uuid.UUID) (*model.UserInfo, error)
}
