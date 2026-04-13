package auth

import (
	"errors"
	"strings"

	"tic-tac-toe/internal/api/http/dto"
	"tic-tac-toe/internal/domain/model"
	"tic-tac-toe/internal/domain/service"

	"github.com/google/uuid"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	userService service.UserService
	jwtProvider *JwtProvider
}

func NewAuthService(userService service.UserService, jwtProvider *JwtProvider) *AuthService {
	return &AuthService{
		userService: userService,
		jwtProvider: jwtProvider,
	}
}

func (s *AuthService) Register(req *dto.SignUpRequest) error {
	return s.userService.CreateUser(req.Login, req.Password)
}

func (s *AuthService) Authorize(authHeader string) (uuid.UUID, error) {
	token, err := s.validateAuthHeader(authHeader)
	if err != nil {
		return uuid.Nil, err
	}

	userID, err := s.jwtProvider.ValidateAccessToken(token)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func (s *AuthService) Authenticate(req *JwtRequest) (*JwtResponse, error) {
	userID, err := s.userService.Authenticate(req.Login, req.Password)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtProvider.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtProvider.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &JwtResponse{
		Type:         "login",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) UpdateAccessToken(authHeader string) (*JwtResponse, error) {
	token, err := s.validateAuthHeader(authHeader)
	if err != nil {
		return nil, err
	}

	userID, err := s.jwtProvider.ValidateRefreshToken(token)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtProvider.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	return &JwtResponse{
		Type:        "update access",
		AccessToken: accessToken,
	}, nil
}

func (s *AuthService) UpdateRefreshToken(authHeader string) (*JwtResponse, error) {
	token, err := s.validateAuthHeader(authHeader)
	if err != nil {
		return nil, err
	}

	userID, err := s.jwtProvider.ValidateRefreshToken(token)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwtProvider.GenerateAccessToken(userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.jwtProvider.GenerateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	return &JwtResponse{
		Type:         "update refresh",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) GetMyInfo(userID string) (*model.UserInfo, error) {
	return s.userService.GetUserInfo(userID)
}

func (s *AuthService) validateAuthHeader(authHeader string) (string, error) {
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", ErrInvalidCredentials
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	return token, nil
}
