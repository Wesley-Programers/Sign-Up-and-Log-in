package repository

import "context"

type User interface {
	Register(ctx context.Context, name, email, password string) error
}

type LoginUser interface {
	VerifyLogin(name, email, password string) error
}

type ChangeName interface {
	ChangeName(currentName, newName string) error
}

type ChangeEmail interface {
	ChangeEmail(currentEmail, newEmail, confirmNewEmail, password string) error
}

type Request interface {
	Request(email string) (error, int)
}

type ResetPassword interface {
	ResetPassword(currentPassword, newPassword, confirmNewPassword string) (error, string)
}

type ValidToken interface {
	ValidToken(token string) error
}

type DeleteAccount interface {
	DeleteAccount(email, password string) error
}
