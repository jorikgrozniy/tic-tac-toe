package middleware

import (
	"net/http"

	"tic-tac-toe/internal/api/http/contextkey"
	"tic-tac-toe/internal/api/http/util"
	"tic-tac-toe/internal/auth"
)

type UserAuthenticator struct {
	authService *auth.AuthService
}

func NewUserAuthenticator(authService *auth.AuthService) *UserAuthenticator {
	return &UserAuthenticator{
		authService: authService,
	}
}

func (ua *UserAuthenticator) RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		userID, err := ua.authService.Authenticate(authHeader)
		if err != nil {
			util.SendError(w, "Authentication required", http.StatusUnauthorized)
			return
		}

		ctx := contextkey.WithUserID(r.Context(), userID)
		next(w, r.WithContext(ctx))
	}
}
