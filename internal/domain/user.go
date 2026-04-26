package domain

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id int
	Name string
	Email string
	PasswordHash string
}

func (user *User) ID() int { return user.Id }
func (user *User) NAME() string { return user.Name }
func (user *User) EMAIL() string { return user.Email }


func (user *User) ChangeEmail(currentEmail, newEmail, confirmNewEmail string) error {
	
	if newEmail != confirmNewEmail {
		return ErrEmailMismatch
	}

	if !strings.EqualFold(currentEmail, user.Email) {
		return ErrEmailMismatch
	}

	if strings.EqualFold(newEmail, user.Email) {
		return ErrEmailIsTheSame
	}

	if !strings.HasSuffix(newEmail, "@example.com") {
		return ErrInvalidEmailFormat
	}

	user.Email = strings.ToLower(newEmail)
	return nil
}


func (user *User) ChangeName(currentName, newName string) error {

	if !strings.EqualFold(currentName, user.Name) {
		return ErrUserNotFound
	}

	if strings.EqualFold(newName, user.Name) {
		return ErrNameIsTheSame
	}

	user.Name = strings.ToLower(newName)
	return nil
}


func (user *User) Login(nameOrEmail, password string) error {

	if strings.HasSuffix(nameOrEmail, "@example.com") {
		if !strings.EqualFold(nameOrEmail, user.Email) {
			return ErrUserNotFound
		}

	} else {
		if !strings.EqualFold(nameOrEmail, user.Name) {
			return ErrUserNotFound
		}
	}

	return nil
}


func (user *User) Register(name, email, password string) (*User, error) {
	if name == "" || email == "" {
		return nil, ErrInvalidCredentials
	}

	cleanName := strings.TrimSpace(name)
	cleanEmail := strings.TrimSpace(email)

	if !strings.HasSuffix(cleanEmail, "@example.com") {
		return nil, ErrInvalidEmailFormat
	}

	return &User{
		Name: strings.ToLower(cleanName),
		Email : strings.ToLower(cleanEmail),
		PasswordHash: password,
	}, nil
}


func (user *User) RequestPasswordReset(email string) error {
	cleanEmail := strings.TrimSpace(email)
	if !strings.HasSuffix(cleanEmail, "@example.com") || cleanEmail == "" {
		return ErrInvalidEmailFormat
	}

	if !strings.EqualFold(cleanEmail, user.Email) {
		return ErrUserNotFound
	}
	
	return nil
}


func (user *User) PasswordValid(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
}


func RestoreEmail(id int, email, passwordHash string) *User {
	return &User{
		Id: id,
		Email: email,
		PasswordHash: passwordHash,
	}
}

func RestoreName(id int, name string) *User {
	return &User{
		Id: id,
		Name: name,
	}
}

func RestoreLogin(id int, name, email, passwordHash string) *User {
	return &User{
		Id: id,
		Name: name,
		Email: email,
		PasswordHash: passwordHash,
	}
}

func RestoreRequest(id int, email string) *User {
	return &User{
		Id: id,
		Email: email,
	}
}