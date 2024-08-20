package authhandler

import (
	"VK-test/core"
	"github.com/stretchr/testify/mock"
)

type MockServiceAuth struct {
	mock.Mock
}

func (m *MockServiceAuth) Login(user core.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}
