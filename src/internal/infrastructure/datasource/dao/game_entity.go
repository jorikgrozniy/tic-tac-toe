package dao

import (
	"time"

	"github.com/google/uuid"
)

type GameEntity struct {
	ID        uuid.UUID    `db:"id"`
	Board     [9]int       `db:"board"`
	Status    EntityStatus `db:"status"`
	Players   []GamePlayer `db:"players"`
	Type      EntityType   `db:"type"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
}

type GamePlayer struct {
	ID   *uuid.UUID `json:"id"`
	IsAI bool       `json:"is_ai"`
	Flag int        `json:"flag"`
}

type EntityStatus string

const (
	StatusWaiting EntityStatus = "waiting"
	StatusTurnX   EntityStatus = "x turn"
	StatusTurnO   EntityStatus = "o turn"
	StatusWinX    EntityStatus = "x won"
	StatusWinO    EntityStatus = "o won"
	StatusDraw    EntityStatus = "draw"
)

type EntityType string

const (
	TypePVP EntityType = "pvp"
	TypePVE EntityType = "pve"
)
