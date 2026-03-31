package service

import (
	"tic-tac-toe/internal/domain/model"

	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(username, password string) error
	Authenticate(username, password string) (uuid.UUID, error)
	GetUserInfo(userID string) (*model.UserInfo, error)
}
