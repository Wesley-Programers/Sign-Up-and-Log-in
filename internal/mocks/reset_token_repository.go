package mocks

import (
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type ResetTokenRepositoryMock struct {
	mock.Mock
}

func (m *ResetTokenRepositoryMock) Save(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	args := m.Called(ctx, userID, tokenHash, expiresAt)
	return args.Error(0)
}

func (m *ResetTokenRepositoryMock) FindValid(ctx context.Context, tokenHash string) (string, error) {
	args := m.Called(ctx, tokenHash)
	return args.String(0), args.Error(1)
}

func (m *ResetTokenRepositoryMock) MarkUsed(ctx context.Context, tokenHash string) error {
	args := m.Called(ctx, tokenHash)
	return args.Error(0)
}

func (m *ResetTokenRepositoryMock) InvalidateAll(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}