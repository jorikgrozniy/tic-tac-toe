package di

import (
	"context"
	"log"
	"net/http"
	"time"

	"tic-tac-toe/config"
	"tic-tac-toe/internal/api/http/handler"
	"tic-tac-toe/internal/api/http/middleware"
	"tic-tac-toe/internal/api/http/router"
	"tic-tac-toe/internal/auth"
	"tic-tac-toe/internal/domain/service"
	"tic-tac-toe/internal/infrastructure/database"
	"tic-tac-toe/internal/infrastructure/datasource/repository"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func Container() fx.Option {
	return fx.Module("app",
		fx.Provide(
			config.New,
		),

		fx.Provide(
			func(lc fx.Lifecycle, config *config.Config) (*pgxpool.Pool, error) {
				return newDatabasePool(lc, config.Database.URL)
			},
		),

		fx.Provide(
			repository.NewUserRepository,
			repository.NewGameRepository,
		),

		fx.Provide(
			service.NewUserService,
			service.NewGameService,
		),

		fx.Provide(
			auth.NewAuthService,
		),

		fx.Provide(
			middleware.NewUserAuthenticator,
		),

		fx.Provide(
			handler.NewAuthHandler,
			handler.NewUserHandler,
			handler.NewGameHandler,
		),

		fx.Provide(
			router.NewRouter,
			routerMuxProvider,
		),

		fx.Provide(func(cfg *config.Config) string {
			return cfg.Server.Port
		}),
	)
}

func newDatabasePool(lc fx.Lifecycle, dbURL string) (*pgxpool.Pool, error) {
	pool, err := database.NewPool(dbURL)
	if err != nil {
		return nil, err
	}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			log.Println("Closing database connection")

			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			pool.Close()

			log.Println("Database connection closed")
			return nil
		},
	})

	return pool, nil
}

func routerMuxProvider(r *router.Router) *http.ServeMux {
	return r.SetupRoutes()
}
