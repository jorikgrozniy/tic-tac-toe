package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/repository"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/dao"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/mapper"
)

var (
	ErrGameNotFound = errors.New("game not found")
	ErrNilGame      = errors.New("game is nil")
	ErrGameInternal = errors.New("internal server error")
)

type gameRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewGameRepository(db *pgxpool.Pool) repository.GameRepository {
	return &gameRepositoryImpl{
		db: db,
	}
}

func (r *gameRepositoryImpl) Save(game *model.CurrentGame) error {
	if game == nil {
		return ErrNilGame
	}

	entity, err := mapper.ToDatasourceGame(game)
	if err != nil {
		return err
	}

	playersJSON, err := json.Marshal(entity.Players)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO games (
			id, board, status, players, type
		) VALUES (
			$1, $2, $3, $4, $5
		)
		ON CONFLICT (id) DO UPDATE SET
			board = $2,
			status = $3,
			players = $4
	`

	if _, err := r.db.Exec(context.Background(), query,
		entity.ID, entity.Board, entity.Status,
		playersJSON, entity.Type,
	); err != nil {
		return ErrGameInternal
	}

	return nil
}

func (r *gameRepositoryImpl) Retrieve(id uuid.UUID) (*model.CurrentGame, error) {
	var game dao.GameEntity
	query := `
		SELECT
			id, board, status, players, type, created_at, updated_at
		FROM games
		WHERE id = $1
	`

	var playersJSON []byte
	if err := r.db.QueryRow(context.Background(), query, id).Scan(
		&game.ID, &game.Board, &game.Status, &playersJSON, &game.Type, &game.CreatedAt, &game.UpdatedAt,
	); err != nil {
		return nil, ErrGameNotFound
	}

	if err := json.Unmarshal(playersJSON, &game.Players); err != nil {
		return nil, err
	}

	domainGame, err := mapper.ToDomainGame(&game)
	if err != nil {
		return nil, err
	}
	return domainGame, nil
}

func (r *gameRepositoryImpl) FindAvailableGames() ([]model.GameInfo, error) {
	var games []dao.GameEntity
	query := `
		SELECT
			id, created_at
		FROM games
		WHERE status = 'waiting' AND type = 'pvp'
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, ErrGameInternal
	}
	defer rows.Close()

	for rows.Next() {
		var game dao.GameEntity
		if err := rows.Scan(&game.ID, &game.CreatedAt); err != nil {
			return nil, ErrGameNotFound
		}
		games = append(games, game)
	}

	return mapper.ToDomainAvailableGames(games), nil
}
