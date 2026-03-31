package dto

type UserInfo struct {
	ID        string `json:"user_id"`
	Username  string `json:"username"`
	UserSince string `json:"user_since"`
}
