package dto

import "github.com/google/uuid"

type MakeMoveRequest struct {
	Board [3][3]int `json:"board" validate:"required"`
}

type CreateGameRequest struct {
	Type string `json:"type" validate:"required"`
}

type GameResponse struct {
	GameID  string       `json:"game_id"`
	Players []GamePlayer `json:"players"`
	Status  string       `json:"status"`
	Board   [3][3]int    `json:"board"`
}

type GamePlayer struct {
	ID   string `json:"id,omitempty"`
	IsAI bool   `json:"is_ai,omitempty"`
	Flag int    `json:"flag,omitempty"`
}

const (
	StatusWaiting = "waiting"
	StatusTurnX   = "turn x"
	StatusTurnO   = "turn o"
	StatusWinX    = "win x"
	StatusWinO    = "win o"
	StatusDraw    = "draw"
)

type AvailableGamesResponse struct {
	Games []GameInfo `json:"games"`
	Total int        `json:"total"`
}

type GameInfo struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt string    `json:"created_at"`
}
