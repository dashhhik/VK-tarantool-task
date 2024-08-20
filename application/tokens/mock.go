package tokens

import (
	"VK-test/core"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type MockPayloadService struct {
	mock.Mock
}

func (m *MockPayloadService) GenerateToken(userPayload core.User) (string, error) {
	args := m.Called(userPayload)
	return args.String(0), args.Error(1)
}

func (m *MockPayloadService) GetUserFromToken(token *jwt.Token) (*core.User, error) {
	args := m.Called(token)
	return args.Get(0).(*core.User), args.Error(1)
}
