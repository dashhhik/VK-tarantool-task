package user_data

import (
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}
