package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"ShieldAuth-API/internal/domain"
)

type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*domain.User), args.Error(1)
}