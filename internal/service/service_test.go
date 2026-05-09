package service

import (
	"context"
	"testing"
	
	"ShieldAuth-API/internal/domain"
	"ShieldAuth-API/internal/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRequestReset(t *testing.T) {
	userRepo := new(mocks.UserRepositoryMock)
	tokenRepo := new(mocks.ResetTokenRepositoryMock)
	securityMock := new(mocks.SecurityMock)

	s := NewService(userRepo, tokenRepo, securityMock)

	userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(&domain.User{Id: 123}, nil)

	securityMock.On("GenerateToken").Return("token123", nil)
	securityMock.On("HashToken", "token123").Return("hashed")

	tokenRepo.On("InvalidateAll", mock.Anything, 123).Return(nil)
	tokenRepo.On("Save", mock.Anything, 123, "hashed", mock.Anything).Return(nil)

	token, err := s.RequestReset(context.Background(), "test@example.com")

	assert.NoError(t, err)
	assert.Equal(t, "token123", token)

	userRepo.AssertExpectations(t)
	tokenRepo.AssertExpectations(t)
	securityMock.AssertExpectations(t)
}