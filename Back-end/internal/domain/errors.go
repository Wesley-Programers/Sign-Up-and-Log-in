package domain

import "errors"

var (
	ErrUserNotFound = errors.New("User not found")
	ErrEmailAlreadyExist = errors.New("Email already exist")
	ErrInvaliData = errors.New("Invalid input")
	ErrInternal = errors.New("Internal error occurred")
)