package repository

import (
	"context"

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

type ResetPassword interface {
	UpdatePassword(ctx context.Context, userID string, passwordHash string) error
}

type DeleteAccount interface {
	DeleteAccount(ctx context.Context, email, password string) error
}
