package mapper

import (
	"tic-tac-toe/internal/domain/model"
	"tic-tac-toe/internal/infrastructure/datasource/dao"
)

func ToDomainGame(entity *dao.GameEntity) (*model.CurrentGame, error) {
	if entity == nil {
		return nil, nil
	}

	players := toDomainPlayers(entity.Players)

	return &model.CurrentGame{
		ID:      entity.ID,
		Board:   model.GameBoard(convert1Dto2D(entity.Board)),
		Status:  toDomainStatus(entity.Status, players),
		Players: players,
		Type:    toDomainType(entity.Type),
	}, nil
}

func ToDatasourceGame(entity *model.CurrentGame) (*dao.GameEntity, error) {
	if entity == nil {
		return nil, nil
	}

	return &dao.GameEntity{
		ID:      entity.ID,
		Board:   convert2Dto1D(entity.Board),
		Status:  toDatasourceStatus(entity),
		Players: toDatasourcePlayers(entity.Players),
		Type:    toDatasourceType(entity.Type),
	}, nil
}

func toDatasourcePlayers(players model.GamePlayers) []dao.GamePlayer {
	return []dao.GamePlayer{
		dao.GamePlayer(players.P1),
		dao.GamePlayer(players.P2),
	}
}

func toDomainPlayers(players []dao.GamePlayer) model.GamePlayers {
	return model.GamePlayers{
		P1: model.GamePlayer(players[0]),
		P2: model.GamePlayer(players[1]),
	}
}

func ToDomainAvailableGames(games []dao.GameEntity) []model.GameInfo {
	domainGames := make([]model.GameInfo, len(games))
	for i, game := range games {
		domainGames[i] = model.GameInfo{
			ID:        game.ID,
			CreatedAt: game.CreatedAt.String(),
		}
	}
	return domainGames
}

func toDomainStatus(status dao.EntityStatus, players model.GamePlayers) model.GameStatus {
	switch status {
	case dao.StatusWaiting:
		return model.GameStatus{
			Type: model.StatusWaiting,
		}
	case dao.StatusTurnX:
		return model.GameStatus{
			Type:   model.StatusTurn,
			Player: players.GetX(),
		}
	case dao.StatusTurnO:
		return model.GameStatus{
			Type:   model.StatusTurn,
			Player: players.GetO(),
		}
	case dao.StatusWinX:
		return model.GameStatus{
			Type:   model.StatusWin,
			Player: players.GetX(),
		}
	case dao.StatusWinO:
		return model.GameStatus{
			Type:   model.StatusWin,
			Player: players.GetO(),
		}
	case dao.StatusDraw:
		return model.GameStatus{
			Type: model.StatusDraw,
		}
	}
	return model.GameStatus{}
}

func toDatasourceStatus(game *model.CurrentGame) dao.EntityStatus {
	switch game.Status.Type {
	case model.StatusWaiting:
		return dao.StatusWaiting
	case model.StatusTurn:
		switch game.Status.Player {
		case game.Players.GetX():
			return dao.StatusTurnX
		case game.Players.GetO():
			return dao.StatusTurnO
		}
	case model.StatusWin:
		switch game.Status.Player {
		case game.Players.GetX():
			return dao.StatusWinX
		case game.Players.GetO():
			return dao.StatusWinO
		}
	case model.StatusDraw:
		return dao.StatusDraw
	}
	return ""
}

func toDomainType(t dao.EntityType) model.GameType {
	switch t {
	case dao.TypePVP:
		return model.TypePVP
	case dao.TypePVE:
		return model.TypePVE
	}
	return ""
}

func toDatasourceType(t model.GameType) dao.EntityType {
	switch t {
	case model.TypePVP:
		return dao.TypePVP
	case model.TypePVE:
		return dao.TypePVE
	}
	return ""
}

func convert2Dto1D(board [3][3]int) [9]int {
	var oneD [9]int
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			oneD[i*3+j] = board[i][j]
		}
	}
	return oneD
}

func convert1Dto2D(board [9]int) [3][3]int {
	var twoD [3][3]int
	for i := 0; i < 9; i++ {
		twoD[i/3][i%3] = board[i]
	}
	return twoD
}
