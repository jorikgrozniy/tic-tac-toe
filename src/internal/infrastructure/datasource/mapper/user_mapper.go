package mapper

import (
	"tic-tac-toe/internal/domain/model"
	"tic-tac-toe/internal/infrastructure/datasource/dao"
)

func ToDomainUser(entity *dao.UserEntity) *model.User {
	if entity == nil {
		return nil
	}

	return &model.User{
		ID:       entity.ID,
		Username: entity.Username,
		Password: entity.Password,
	}
}

func ToDomainUserInfo(entity *dao.UserEntity) *model.UserInfo {
	if entity == nil {
		return nil
	}

	return &model.UserInfo{
		ID:        entity.ID,
		Username:  entity.Username,
		UserSince: entity.CreatedAt.String(),
	}
}

func ToDatasourceUser(entity *model.User) *dao.UserEntity {
	if entity == nil {
		return nil
	}

	return &dao.UserEntity{
		ID:       entity.ID,
		Username: entity.Username,
		Password: entity.Password,
	}
}
