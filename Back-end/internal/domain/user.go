package domain

import (
	"strings"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID int
	Name string
	Email string
	PasswordHash string
}

func (user *User) ChangeEmail(currentEmail, newEmail, confirmEmail string) error {
	currentEmail = strings.ToLower(strings.TrimSpace(currentEmail))
	newEmail = strings.ToLower(strings.TrimSpace(newEmail))
	confirmEmail = strings.ToLower(strings.TrimSpace(confirmEmail))

	if newEmail == user.Email {
		return ErrEmailIsTheSame
	}

	if newEmail != confirmEmail {
		return ErrEmailMismatch
	}

	if newEmail == "" {
		return ErrInvalidData
	}

	user.Email = newEmail
	return nil
}

func (user *User) PasswordValid(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}