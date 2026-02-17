package service

import (
	"encoding/hex"
	"errors"
	"time"
	"unicode"

	"crypto/rand"

	"index/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Repository repository.User
}

type VerifyLogin struct {
	Repository repository.LoginUser
}

type ChangeName struct {
	Repository repository.ChangeName
}

type ChangeEmail struct {
	Repository repository.ChangeEmail
}

type Request struct {
	Repository repository.Request
}

type ResetPassword struct {
	Repository repository.ResetPassword
}

type DeleteAccount struct {
	Repository repository.DeleteAccount
}


func (user *User) SaveData(name, email, password string) error {
	validPassword, message := VerifyPassword(password)
	if !validPassword {
		return errors.New(message)
	}

	hash, err := HashPassword(password)
	if err != nil && validPassword {
		return err
	}

	return user.Repository.NewMysqlInsert(name, email, string(hash))
}


func (verifyLogin *VerifyLogin) VerifyLoginFunction(name, email, password string) error {
	return verifyLogin.Repository.NewVerifyLogin(name, email, password)
}


func (changeName *ChangeName) ChangeNameFunction(currentName, newName string) error {
	return changeName.Repository.ChangeName(currentName, newName)
}


func (changeEmail *ChangeEmail) ChangeEmailFunction(currentEmail, newEmail, confirmNewEmail, password string) error {
	return changeEmail.Repository.ChangeEmail(currentEmail, newEmail, confirmNewEmail, password)
}


func (request *Request) RequestFunction(email string) (error, string) {
	err, id := request.Repository.Request(email)
	if err != nil {
		return err, ""
		
	} else {
		token, err := GenerateTokens()
		if err != nil {
			return err, ""
		}
		expiresAt := time.Now().Add(10 *time.Minute)
		err = repository.InsertInto(id, token, expiresAt)
		if err != nil {
			return err, ""
		}

		link := GenerateLink(token)
		return nil, link

	}
}


func (resetPasword *ResetPassword) ResetPasswordFunction(currentPassword, newPassword, confirmNewPassword string) error {
	validPassword, message := VerifyPassword(newPassword)
	if !validPassword {
		return errors.New(message)
	}

	hash, err := HashPassword(newPassword)
	if err != nil && validPassword {
		return err
	}

	err, email := resetPasword.Repository.ResetPassword(currentPassword, newPassword, confirmNewPassword)
	if err != nil {
		_, err = repository.LimitOfAttempts(email)
		if err != nil {
			return err
		}
		return err
	}
	err = repository.UpdatePassword(hash, email)
	if err != nil {
		_, err = repository.LimitOfAttempts(email)
		if err != nil {
			return err
		}
		return err
	}

	return nil
}


func (delete *DeleteAccount) DeleteAccountFunction(email, password string) error {
	return delete.Repository.DeleteAccount(email, password)
}


func VerifyPassword(password string) (bool, string) {

	if len(password) < 8 {
		return false, "short password"
	}

	if len(password) >= 150 {
		return false, "long password"
	}

	var (
		upper = false
		lower = false
		special = false
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			upper = true
		
		case unicode.IsLower(char):
			lower = true

		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			special = true
		}
	}

	if !upper {
		return false, "at least one capital letter"
	}
	if !lower {
		return false, "at least one lowercase letter"
	}
	if !special {
		return false, "at least one special character"
	}

	return true, ""

}


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}


func GenerateTokens() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}


func GenerateLink(token string) string {
	return "http://127.0.0.1:5500/reset?token=" + token
}


func StartToRemoverExpiredTokens() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		repository.RemoveExpiredToken()
	}
}
