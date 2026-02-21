package service

import (
	"context"
	"encoding/hex"
	"errors"
	"time"
	"fmt"
	"unicode"

	"crypto/rand"
	"crypto/sha256"

	"index/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Register struct {
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

type ValidToken struct {
	Repository repository.ValidToken
}

func NewUserStruct(repository repository.User) *Register {
	return &Register{
		Repository: repository,
	}
}

var smallTest = make([]string, 0, 1)

func (register *Register) RegisterFunction(ctx context.Context, name, email, password string) error {
	validPassword, message := VerifyPassword(password)
	if !validPassword {
		return errors.New(message)
	}

	hash, err := HashPassword(password)
	if err != nil && validPassword {
		return err
	}

	return register.Repository.Register(ctx, name, email, string(hash))
}


func (verifyLogin *VerifyLogin) VerifyLoginFunction(name, email, password string) error {
	return verifyLogin.Repository.VerifyLogin(name, email, password)
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

		} else {
			smallTest = append(smallTest, token)

			ctx, cancel := context.WithCancel(context.Background())
			go func (ctx context.Context) {
				select {
				case <- time.After(15 *time.Second):
					smallTest = smallTest[:0]

				case <-ctx.Done():
					fmt.Println("CLEAN")
					return
				}
			}(ctx)
			cancel()
		}	

		expiresAt := time.Now().Add(10 * time.Minute)
		tokenHash := tokenHash(token)
		err = repository.InsertInto(id, tokenHash, expiresAt)
		if err != nil {
			return err, ""
		}

		return nil, token

	}
}


func (validToken *ValidToken) ValidTokenFunction() error {
	token := smallTest[0]
	var test string
	tokenHash := tokenHash(token)
	err := repository.Test(tokenHash, test)
	if err != nil {
		return err
	} 

	if tokenHash != test {
		return errors.New("ERROR")
	}
	smallTest = smallTest[:0]
	return nil
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


func StartToRemoverExpiredTokens() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		repository.RemoveExpiredToken()
	}
}


func tokenHash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
