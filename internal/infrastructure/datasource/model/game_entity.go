package model

import (
	"time"

	"github.com/google/uuid"
)

type GameEntity struct {
	ID        uuid.UUID
	Board     [3][3]int
	CreatedAt time.Time
	UpdatedAt time.Time
}
