package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type LimiterMock struct {
	mock.Mock
}

func (m *LimiterMock) CheckLimit(ctx context.Context, key string, maxAttempts int, window time.Duration) (bool, error) {
	args := m.Called(ctx, key, maxAttempts, window)
	return args.Bool(0), args.Error(1)
}