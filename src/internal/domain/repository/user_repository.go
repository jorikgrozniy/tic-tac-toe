package repository

import (
	"tic-tac-toe/internal/domain/model"

	"github.com/google/uuid"
)

type UserRepository interface {
	Save(game *model.User) error
	FindByID(id uuid.UUID) (*model.User, error)
	FindByUsername(username string) (*model.User, error)
	FindUserInfo(id uuid.UUID) (*model.UserInfo, error)
}
