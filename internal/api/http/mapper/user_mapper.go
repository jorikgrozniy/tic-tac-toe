package mapper

import (
	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/dto"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
)

func ToDTOUserInfo(userInfo *model.UserInfo) dto.UserInfo {
	if userInfo == nil {
		return dto.UserInfo{}
	}

	return dto.UserInfo{
		ID:        userInfo.ID.String(),
		Username:  userInfo.Username,
		UserSince: userInfo.UserSince,
	}
}
