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
	userRepo   *mocks.UserRepositoryMock
	security   *mocks.SecurityMock
	resetStore *mocks.ResetStoreMock
	limiter    *mocks.LimiterMock
	service    Service
}

func newDeps() *deps {
	userRepo := new(mocks.UserRepositoryMock)
	security := new(mocks.SecurityMock)
	resetStore := new(mocks.ResetStoreMock)
	limiter := new(mocks.LimiterMock)

	resetStore.On("Save", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	limiter.On("CheckLimit", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(true, nil)

	return &deps{
		userRepo:   userRepo,
		security:   security,
		resetStore: resetStore,
		limiter:    limiter,
		service:    NewService(userRepo, security, resetStore, limiter),
	}
}

func TestRequestReset(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(d *deps)
		wantToken string
		wantErr   error
	}{
		{
			name: "success",
			setup: func(d *deps) {
				d.userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(&domain.User{Id: 123}, nil)

				d.security.On("GenerateToken").Return("token123", nil)

				d.resetStore.On("Save", mock.Anything, "token123", 123, mock.Anything).Return(nil)
			},
			wantToken: "token123",
			wantErr:   nil,
		},

		{
			name: "user not found",
			setup: func(d *deps) {
				d.userRepo.On("GetByEmail", mock.Anything, mock.Anything).Return(nil, domain.ErrUserNotFound)
			},
			wantToken: "",
			wantErr:   nil,
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
			d.security.AssertExpectations(t)
		})
	}
}
