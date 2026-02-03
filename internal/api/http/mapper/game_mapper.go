package mapper

import (
	"github.com/google/uuid"
	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/dto"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
)

func ToDTOGame(game *model.CurrentGame) dto.GameResponse {
	if game == nil {
		return dto.GameResponse{}
	}

	return dto.GameResponse{
		GameID:  game.ID.String(),
		Players: toDTOPlayers(game.Players),
		Status:  toDTOStatus(game),
		Board:   game.Board,
	}
}

func toDTOPlayers(players model.GamePlayers) []dto.GamePlayer {
	var id1, id2 string

	if players.P1.ID != nil {
		id1 = players.P1.ID.String()
	}

	if players.P2.ID != nil {
		id2 = players.P2.ID.String()
	}

	return []dto.GamePlayer{
		{
			ID:   id1,
			IsAI: players.P1.IsAI,
			Flag: players.P1.Flag,
		},
		{
			ID:   id2,
			IsAI: players.P2.IsAI,
			Flag: players.P2.Flag,
		},
	}
}

func toDTOStatus(game *model.CurrentGame) string {
	switch game.Status.Type {
	case model.StatusWaiting:
		return dto.StatusWaiting
	case model.StatusTurn:
		switch *game.Status.Player {
		case *game.Players.GetX():
			return dto.StatusTurnX
		case *game.Players.GetO():
			return dto.StatusTurnO
		}
	case model.StatusWin:
		switch *game.Status.Player {
		case *game.Players.GetX():
			return dto.StatusWinX
		case *game.Players.GetO():
			return dto.StatusWinO
		}
	case model.StatusDraw:
		return dto.StatusDraw
	}
	return ""
}

func FromDTOGame(gameID string, req dto.MakeMoveRequest) *model.CurrentGame {
	id, err := uuid.Parse(gameID)
	if err != nil {
		return nil
	}

	return &model.CurrentGame{
		ID:    id,
		Board: model.GameBoard(req.Board),
	}
}

func ToDTOAvailableGames(games []model.GameInfo, total int) dto.AvailableGamesResponse {
	dtoGames := make([]dto.GameInfo, len(games))

	for i, game := range games {
		dtoGames[i] = dto.GameInfo{
			ID:        game.ID,
			CreatedAt: game.CreatedAt,
		}
	}

	return dto.AvailableGamesResponse{
		Games: dtoGames,
		Total: total,
	}
}
