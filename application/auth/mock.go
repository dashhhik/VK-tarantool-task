package auth

import (
	"VK-test/core"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type MockTokenService struct {
	mock.Mock
}

func (m *MockTokenService) GenerateToken(userPayload core.User) (string, error) {
	args := m.Called(userPayload)
	return args.String(0), args.Error(1)
}

func (m *MockTokenService) GetUserFromToken(token *jwt.Token) (*core.User, error) {
	args := m.Called(token)
	return args.Get(0).(*core.User), args.Error(1)
}

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Get(username string) (string, error) {
	args := m.Called(username)
	return args.String(0), args.Error(1)
}
