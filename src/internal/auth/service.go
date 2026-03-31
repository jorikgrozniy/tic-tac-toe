package auth

import (
	"encoding/base64"
	"errors"
	"strings"

	"tic-tac-toe/internal/api/http/dto"
	"tic-tac-toe/internal/domain/service"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	userService service.UserService
}

func NewAuthService(userService service.UserService) *AuthService {
	return &AuthService{
		userService: userService,
	}
}

func (s *AuthService) Register(req *dto.SignUpRequest) error {
	return s.userService.CreateUser(req.Login, req.Password)
}

func (s *AuthService) Authenticate(authHeader string) (uuid.UUID, error) {
	if !strings.HasPrefix(authHeader, "Basic ") {
		return uuid.Nil, ErrInvalidCredentials
	}

	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	credentials := string(decoded)
	parts := strings.SplitN(credentials, ":", 2)
	if len(parts) != 2 {
		return uuid.Nil, ErrInvalidCredentials
	}

	username := parts[0]
	password := parts[1]
	userID, err := s.userService.Authenticate(username, password)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}
