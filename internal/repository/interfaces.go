package repository

import (
	"context"
	"time"

	"ShieldAuth-API/internal/domain"
)

type User interface {
	Register(ctx context.Context, a, b, c string) error
}

type LoginUser interface {
	GetByCredentials(ctx context.Context, u domain.User) (*domain.User, error)
}

type ChangeName interface {
	GetID(ctx context.Context, id int) (*domain.User, error)
	UpdateName(ctx context.Context, user *domain.User) error
}

type ChangeEmail interface {
	GetID(ctx context.Context, id int) (*domain.User, error)
	UpdateEmail(ctx context.Context, user *domain.User) error
}

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

type ResetTokenRepository interface {
	Save(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error
	FindValid(ctx context.Context, tokenHash string) (string, error)
	MarkUsed(ctx context.Context, tokenHash string) error
	InvalidateAll(ctx context.Context, userID int) error
}

type ResetPassword interface {
	ResetPassword(ctx context.Context, currentPassword, newPassword, confirmNewPassword string) (error, string)
}

type DeleteAccount interface {
	DeleteAccount(ctx context.Context, email, password string) error
}
