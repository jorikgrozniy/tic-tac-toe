package main

import (
	"log"
	"net/http"

	"go.uber.org/fx"

	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/handlers"
	"github.com/jorikgrozniy/tic-tac-toe/internal/application"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/repository"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/storage"
)

func main() {
	app := fx.New(
		fx.Provide(
			storage.NewGameStorage,
			repository.NewGameRepositoryImpl,
			application.NewGameAppService,
			handlers.NewGameHandler,
			createMux,
		),

		fx.Invoke(startServer),
	)

	app.Run()
}

func createMux(handler *handlers.GameHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /game/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handler.GetGameHandler(w, r, id)
	})

	mux.HandleFunc("POST /game/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		handler.MakeMoveHandler(w, r, id)
	})

	mux.HandleFunc("POST /game", handler.CreateGameHandler)
	return mux
}

func applyMiddleware(mux *http.ServeMux) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.Method, r.URL.Path)

		w.Header().Set("Content-Type", "application/json")

		mux.ServeHTTP(w, r)
	})
}

func startServer(mux *http.ServeMux) {
	port := ":8080"
	log.Printf("Server starting on http://localhost%s", port)
	log.Println("Available endpoints:")
	log.Println("	GET /game/{id}          - Get game state")
	log.Println("	POST /game/{id}   - Make a move")
	log.Println("	POST /game              - Create new game")

	if err := http.ListenAndServe(":8080", applyMiddleware(mux)); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
