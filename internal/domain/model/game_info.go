package model

import "github.com/google/uuid"

type GameInfo struct {
	ID        uuid.UUID
	CreatedAt string
}
