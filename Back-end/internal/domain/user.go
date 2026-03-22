package domain

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id int
	name string
	email string
	PasswordHash string
}

func (user *User) ID() int { return user.id }
func (user *User) Name() string { return user.name }
func (user *User) Email() string { return user.email }

func (user *User) ChangeEmail(currentEmail, newEmail, confirmNewEmail string) error {
	
	if !strings.EqualFold(currentEmail, user.email) {
		return ErrEmailMismatch
	}

	if strings.EqualFold(newEmail, user.email) {
		return ErrEmailIsTheSame
	}

	user.email = strings.ToLower(newEmail)
	return nil
}

func (user *User) ChangeName(currentName, newName, confirmNewEmail string) error {

	if !strings.EqualFold(currentName, user.name) {
		return ErrUserNotFound
	}

	if strings.EqualFold(newName, user.name) {
		return ErrNameIsTheSame
	}

	user.name = strings.ToLower(newName)
	return nil
}

func (user *User) PasswordValid(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) == nil
}

func Restore(id int, email, passwordHash string) *User {
	return &User{
		id: id,
		email: email,
		PasswordHash: passwordHash,
	}
}

func RestoreName(id int, name string) *User {
	return &User{
		id: id,
		name: name,
	}
}