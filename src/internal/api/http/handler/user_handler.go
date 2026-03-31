package handler

import (
	"encoding/json"
	"net/http"

	"tic-tac-toe/internal/api/http/mapper"
	"tic-tac-toe/internal/api/http/util"
	"tic-tac-toe/internal/domain/service"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

func (h *UserHandler) GetUserInfoHandler(w http.ResponseWriter, r *http.Request, userID string) {
	if userID == "" {
		util.SendError(w, "user ID is required", http.StatusBadRequest)
		return
	}

	userInfo, err := h.userService.GetUserInfo(userID)
	if err != nil {
		util.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := mapper.ToDTOUserInfo(userInfo)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		util.SendError(w, "error creating response", http.StatusInternalServerError)
		return
	}
}
