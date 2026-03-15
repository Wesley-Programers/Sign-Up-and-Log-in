package repository

import (
	"context"
	"index/Back-end/internal/domain"
)

type User interface {
	Register(ctx context.Context, name, email, password string) error
}

type LoginUser interface {
	VerifyLogin(ctx context.Context, name, email, password string) (error, string, int)
}

type ChangeName interface {
	ChangeName(ctx context.Context, currentName, newName string) error
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
