package auth

import (
	"VK-test/core"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateToken(userPayload core.User) (string, error)
	GetUserFromToken(token *jwt.Token) (*core.User, error)
}

type UserRepo interface {
	Get(key string) (string, error)
}

type Service struct {
	TokenService TokenService
	UserRepo     UserRepo
}

func NewService(tokenService TokenService, userRepo UserRepo) *Service {
	return &Service{TokenService: tokenService, UserRepo: userRepo}
}

func (s *Service) Login(user core.User) (string, error) {
	if user.Username == "" {
		return "", core.NewCustomError(400, "invalid request format")
	}

	passwordCheck, err := s.UserRepo.Get(user.Username)
	if err != nil {
		if errors.Is(err, core.ErrKeyNotFound) {
			return "", core.NewCustomError(404, "user not found")
		}
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if passwordCheck != user.Password {
		return "", core.NewCustomError(401, "invalid credentials")
	}

	token, err := s.TokenService.GenerateToken(user)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}
