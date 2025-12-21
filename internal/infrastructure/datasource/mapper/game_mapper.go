package mapper

import (
	domainModel "github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
	dsModel "github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/model"
)

func ToDomainEntity(entity *dsModel.GameEntity) *domainModel.CurrentGame {
	if entity == nil {
		return nil
	}

	return &domainModel.CurrentGame{
		ID:    entity.ID,
		Board: domainModel.GameBoard(entity.Board),
	}
}

func ToDataSourceEntity(entity *domainModel.CurrentGame) *dsModel.GameEntity {
	if entity == nil {
		return nil
	}

	return &dsModel.GameEntity{
		ID:    entity.ID,
		Board: [3][3]int(entity.Board),
	}
}
