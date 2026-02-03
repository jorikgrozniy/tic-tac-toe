package model

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID
	Username string
	Password string
}

const (
	MinUsernameLength = 3
	MaxUsernameLength = 50
	MinPasswordLength = 8
	MaxPasswordLength = 255
)
