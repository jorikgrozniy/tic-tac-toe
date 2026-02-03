package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/model"
	"github.com/jorikgrozniy/tic-tac-toe/internal/domain/repository"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/dao"
	"github.com/jorikgrozniy/tic-tac-toe/internal/infrastructure/datasource/mapper"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrNilUser      = errors.New("user is nil")
	ErrUserInternal = errors.New("internal server error")
)

type userRepositoryImpl struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) repository.UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) Save(user *model.User) error {
	if user == nil {
		return ErrNilUser
	}

	entity := mapper.ToDatasourceUser(user)

	query := `
		INSERT INTO users (
			id, username, password
		) VALUES (
			$1, $2, $3
		)
	`

	_, err := r.db.Exec(context.Background(), query, entity.ID, entity.Username, entity.Password)
	if err != nil {
		return ErrGameInternal
	}

	return nil
}

func (r *userRepositoryImpl) loadUser(id uuid.UUID) (*dao.UserEntity, error) {
	var user dao.UserEntity
	query := `
		SELECT
			id, username, password, created_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.CreatedAt,
	)

	if err != nil {
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (r *userRepositoryImpl) FindByID(id uuid.UUID) (*model.User, error) {
	user, err := r.loadUser(id)
	return mapper.ToDomainUser(user), err
}

func (r *userRepositoryImpl) FindUserInfo(id uuid.UUID) (*model.UserInfo, error) {
	user, err := r.loadUser(id)
	return mapper.ToDomainUserInfo(user), err
}

func (r *userRepositoryImpl) FindByUsername(username string) (*model.User, error) {
	var user dao.UserEntity
	query := `
		SELECT
			id, username, password, created_at
		FROM users
		WHERE username = $1
	`

	err := r.db.QueryRow(context.Background(), query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.CreatedAt,
	)

	if err != nil {
		return nil, ErrUserNotFound
	}
	return mapper.ToDomainUser(&user), nil
}
