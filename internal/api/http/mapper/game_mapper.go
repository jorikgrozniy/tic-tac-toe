package mapper

import (
	"github.com/google/uuid"
	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/dto"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
)

func ToDTO(game *model.CurrentGame, status int) dto.Response {
	if game == nil {
		return dto.Response{}
	}

	var statusStr string
	switch status {
	case -1:
		statusStr = "Playing"
	case 0:
		statusStr = "Draw"
	case 1:
		statusStr = "WinX"
	case 2:
		statusStr = "WinO"
	}

	return dto.Response{
		GameID: game.ID.String(),
		Board:  game.Board,
		Status: statusStr,
	}
}

func FromDTO(gameID string, req dto.MakeMoveRequest) *model.CurrentGame {
	id, err := uuid.Parse(gameID)
	if err != nil {
		return nil
	}

	return &model.CurrentGame{
		ID:    id,
		Board: model.GameBoard(req.Board),
	}
}
