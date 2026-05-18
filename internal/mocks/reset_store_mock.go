package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type ResetStoreMock struct {
	mock.Mock
}

func (m *ResetStoreMock) Get(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func (m *ResetStoreMock) Delete(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *ResetStoreMock) Save(ctx context.Context, token string, userID int, ttl time.Duration) error {
	args := m.Called(ctx, token, userID, ttl)
	return args.Error(0)
}
