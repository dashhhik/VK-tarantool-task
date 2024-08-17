package auth

import (
	"VK-test/core"
	"github.com/golang-jwt/jwt/v5"
)

type TokenService interface {
	GenerateToken(userPayload core.User) (string, error)
	GetUserFromToken(token *jwt.Token) (*core.User, error)
}

type UserRepo interface {
	Get(key string) (interface{}, error)
}

type Service struct {
	TokenService TokenService
	UserRepo     UserRepo
}

func NewService(tokenService TokenService, userRepo UserRepo) *Service {
	return &Service{TokenService: tokenService, UserRepo: userRepo}
}

func (s Service) Login(user core.User) (string, error) {
	passwordCheck, err := s.UserRepo.Get(user.Username)
	if err != nil {
		return "", err
	}

	if passwordCheck != user.Password {
		return "", err
	}

	token, err := s.TokenService.GenerateToken(user)
	if err != nil {
		return "", err
	}

	return token, nil

}
