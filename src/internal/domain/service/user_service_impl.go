package service

import (
	"errors"
	"regexp"

	"tic-tac-toe/internal/domain/model"
	"tic-tac-toe/internal/domain/repository"

	"github.com/google/uuid"
)

var (
	ErrInvalidPassword           = errors.New("invalid password")
	ErrUserAlreadyExists         = errors.New("user already exists")
	ErrInvalidUsernameLength     = errors.New("username must be between 3 and 50 characters")
	ErrInvalidUsernameCharacters = errors.New("username can only contain letters, numbers and underscores")
	ErrInvalidPasswordLength     = errors.New("password must be between 8 and 255 characters")
	ErrInvalidPasswordCharacters = errors.New("password must contain both letters and numbers")
)

type userServiceImpl struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userServiceImpl{
		repo: repo,
	}
}

func (s *userServiceImpl) CreateUser(username, password string) error {
	if err := s.validateUsername(username); err != nil {
		return err
	}

	if err := s.validatePassword(password); err != nil {
		return err
	}

	exists, err := s.repo.FindByUsername(username)
	if err == nil && exists != nil {
		return ErrUserAlreadyExists
	}

	user := &model.User{
		ID:       uuid.New(),
		Username: username,
		Password: password,
	}

	if err := s.repo.Save(user); err != nil {
		return err
	}

	return nil
}

func (s *userServiceImpl) Authenticate(username, password string) (uuid.UUID, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return uuid.Nil, err
	}

	if !s.checkPassword(password, user) {
		return uuid.Nil, ErrInvalidPassword
	}

	return user.ID, nil
}

func (s *userServiceImpl) GetUserInfo(userID string) (*model.UserInfo, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	userInfo, err := s.repo.FindUserInfo(parsedID)
	if err != nil {
		return nil, err
	}

	return userInfo, nil
}

func (s *userServiceImpl) checkPassword(password string, user *model.User) bool {
	return user.Password == password
}

func (s *userServiceImpl) validateUsername(username string) error {
	if len(username) < model.MinUsernameLength || len(username) > model.MaxUsernameLength {
		return ErrInvalidUsernameLength
	}

	matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]+$`, username)
	if !matched {
		return ErrInvalidUsernameCharacters
	}

	return nil
}

func (s *userServiceImpl) validatePassword(password string) error {
	if len(password) < model.MinPasswordLength || len(password) > model.MaxPasswordLength {
		return ErrInvalidPasswordLength
	}

	hasNumber, _ := regexp.MatchString(`[0-9]`, password)
	hasLetter, _ := regexp.MatchString(`[a-zA-Z]`, password)
	if !hasLetter || !hasNumber {
		return ErrInvalidPasswordCharacters
	}

	return nil
}
