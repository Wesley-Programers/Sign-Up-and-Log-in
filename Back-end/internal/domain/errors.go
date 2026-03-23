package domain

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email, name or password")
	ErrUserNotFound = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exist")
	ErrInvalidPassword = errors.New("invalid password")

	ErrInvalidData = errors.New("invalid input")
	ErrInvalidEmailFormat = errors.New("invalid email format")

	ErrNameIsTheSame = errors.New("the new name is the same as the current one")
	ErrEmailIsTheSame = errors.New("the new email is the same as the current one")
	ErrEmailMismatch = errors.New("the provided current email does not match our records")
	ErrEmailsDoNotMismatch = errors.New("the new email and its confirmation do not match")

	ErrInternal = errors.New("an internal error occurred")
	ErrRateLimitExceeded = errors.New("too many requests, please try again later")
)