package service

import (
	"context"
	"testing"

	"ShieldAuth-API/internal/domain"
	"ShieldAuth-API/internal/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type deps struct {
	userRepo *mocks.UserRepositoryMock
	tokenRepo *mocks.ResetTokenRepositoryMock
	security *mocks.SecurityMock
	service Service	
}

func newDeps() *deps {
	userRepo := new(mocks.UserRepositoryMock)
	tokenRepo := new(mocks.ResetTokenRepositoryMock)
	security := new(mocks.SecurityMock)


	return &deps {
		userRepo: userRepo,
		tokenRepo: tokenRepo,
		security: security,
		service: NewService(userRepo, tokenRepo, security),
	}
}

func TestRequestReset(t *testing.T) {
	tests := []struct {
		name string
		setup func(d *deps)
		wantToken string
		wantErr error
	}{
		{
			name: "success",
			setup: func(d *deps) {
				d.userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(&domain.User{Id: 123}, nil)
				
				d.security.On("GenerateToken").Return("token123", nil)

				d.security.On("TokenHash", "token123").Return("hashed")

				d.tokenRepo.On("InvalidateAll", mock.Anything, 123).Return(nil)

				d.tokenRepo.On("Save", mock.Anything, 123, "hashed", mock.Anything).Return(nil)

			},
			wantToken: "token123",
			wantErr: nil,
		},

		{
			name: "user not found",
			setup: func(d *deps) {
				d.userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(nil, domain.ErrUserNotFound)
			},
			wantErr: domain.ErrUserNotFound,
		},
	}
	
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			d := newDeps()

			if tt.setup != nil {
				tt.setup(d)
			}

			token, err := d.service.RequestReset(context.Background(), "test@example.com")
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.wantToken, token)
			}

			d.userRepo.AssertExpectations(t)
			d.tokenRepo.AssertExpectations(t)
			d.security.AssertExpectations(t)
		})
	}
}