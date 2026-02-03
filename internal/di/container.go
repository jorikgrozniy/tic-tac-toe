package di

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/handler"
	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/middleware"
	"github.com/jorikgrozniy/tic-tac-toe/internal/api/http/router"
	"github.com/jorikgrozniy/tic-tac-toe/internal/auth"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/service"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/database"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/repository"
	"go.uber.org/fx"
)

func Container() fx.Option {
	return fx.Module("app",
		fx.Provide(
			newDatabasePool,
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
	)
}

func newDatabasePool(lc fx.Lifecycle) (*pgxpool.Pool, error) {
	pool, err := database.NewPool()
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
