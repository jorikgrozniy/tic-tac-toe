package service

import (
	"github.com/google/uuid"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
)

type UserService interface {
	CreateUser(username, password string) error
	Authenticate(username, password string) (uuid.UUID, error)
	GetUserInfo(userID string) (*model.UserInfo, error)
}
