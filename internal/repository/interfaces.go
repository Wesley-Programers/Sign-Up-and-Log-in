package repository

import (
	"context"
	"ShieldAuth-API/internal/domain"
)

type User interface {
	Register(ctx context.Context, name, email, password string) error
}

type LoginUser interface {
	GetByCredentials(ctx context.Context, login string) (*domain.User, error)
}

type ChangeName interface {
	GetID(ctx context.Context, id int) (*domain.User, error)
	UpdateName(ctx context.Context, user *domain.User) error
}

type ChangeEmail interface {
	GetID(ctx context.Context, id int) (*domain.User, error)
	UpdateEmail(ctx context.Context, user *domain.User) error
}

type Request interface {
	Request(ctx context.Context, email string) (error, int)
}

type ResetPassword interface {
	ResetPassword(ctx context.Context, currentPassword, newPassword, confirmNewPassword string) (error, string)
}

type ValidToken interface {
	ValidToken(ctx context.Context, token, secondToken string) error
}

type DeleteAccount interface {
	DeleteAccount(ctx context.Context, email, password string) error
}
