package domain

import "errors"

var (
	ErrUserNotFound = errors.New("User not found")
	ErrEmailAlreadyExist = errors.New("Email already exist")
	ErrInvalidData = errors.New("Invalid input")
	ErrInternal = errors.New("Internal error occurred")
	ErrInvalidPassword = errors.New("Invalid password")
	ErrInvalidCredentials = errors.New("Invalid email, name or password")
	ErrNameIsTheSame = errors.New("The new name is the same as the current one")
	ErrEmailIsTheSame = errors.New("The new email is the same as the current one")
	ErrEmailMismatch = errors.New("The new email and its confirmation are not the same")
)