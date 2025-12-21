package dto

type MakeMoveRequest struct {
	Board [3][3]int `json:"board" validate:"required"`
}

type Response struct {
	GameID string    `json:"game_id"`
	Board  [3][3]int `json:"board"`
	Status string    `json:"status"`
}
