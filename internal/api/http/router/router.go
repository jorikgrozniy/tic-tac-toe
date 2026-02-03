package router

import (
	"net/http"

	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/handler"
	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/middleware"
)

type Router struct {
	authHandler   *handler.AuthHandler
	userHandler   *handler.UserHandler
	gameHandler   *handler.GameHandler
	authenticator *middleware.UserAuthenticator
}

func NewRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	gameHandler *handler.GameHandler,
	authenticator *middleware.UserAuthenticator,
) *Router {
	return &Router{
		authHandler:   authHandler,
		userHandler:   userHandler,
		gameHandler:   gameHandler,
		authenticator: authenticator,
	}
}

func (rt *Router) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	// auth
	mux.HandleFunc("POST /auth/register", rt.authHandler.RegisterHandler)
	mux.HandleFunc("POST /auth/login", rt.authHandler.LoginHandler)

	// user
	mux.HandleFunc("GET /users/{id}", rt.authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		rt.userHandler.GetUserInfoHandler(w, r, id)
	}))

	// game
	mux.HandleFunc("GET /games", rt.authenticator.RequireAuth(rt.gameHandler.GetAvailableGamesHandler))
	mux.HandleFunc("POST /games/create", rt.authenticator.RequireAuth(rt.gameHandler.CreateGameHandler))
	mux.HandleFunc("GET /games/{id}", rt.authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		rt.gameHandler.GetGameHandler(w, r, id)
	}))
	mux.HandleFunc("POST /games/{id}/move", rt.authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		rt.gameHandler.MakeMoveHandler(w, r, id)
	}))
	mux.HandleFunc("POST /games/{id}/join", rt.authenticator.RequireAuth(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		rt.gameHandler.JoinGameHandler(w, r, id)
	}))

	return mux
}
