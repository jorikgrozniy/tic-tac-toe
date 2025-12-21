package model

import "github.com/google/uuid"

type CurrentGame struct {
	ID    uuid.UUID
	Board GameBoard
}
