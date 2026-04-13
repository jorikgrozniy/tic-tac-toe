package main

import (
	"context"
	"log"
	"net/http"

	"go.uber.org/fx"

	"tic-tac-toe/internal/api/http/middleware"
	"tic-tac-toe/internal/di"
)

func main() {
	app := fx.New(
		di.Container(),
		fx.Invoke(startServer),
	)

	app.Run()
}

func startServer(lc fx.Lifecycle, mux *http.ServeMux, port string) {
	server := &http.Server{
		Addr:    port,
		Handler: middleware.Logger(mux),
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("Server starting on %s", server.Addr)
				printInfo()

				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Fatal("Server failed to start:", err)
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Closing server")
			server.Shutdown(ctx)
			log.Println("Server closed")

			return nil
		},
	})
}

func printInfo() {
	log.Println("Available endpoints:")
	log.Println("	POST /auth/register        - Register")
	log.Println("	POST /auth/login           - Login")
	log.Println("	POST /auth/update/access   - Update access token")
	log.Println("	POST /auth/update/refresh  - Update both tokens")
	log.Println("	GET /auth/me               - Get your info")

	log.Println("	GET /users/{id}            - Get user info")

	log.Println("	GET /games                 - List of available games")
	log.Println("	GET /games/{id}            - Get game state")
	log.Println("	POST /games/create         - Create new game")
	log.Println("	POST /games/{id}/move      - Make a move")
	log.Println("	POST /games/{id}/join      - Join game")
}
