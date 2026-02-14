package service

import (
	"encoding/hex"
	"errors"
	"io"
	"os"
	"time"
	"unicode"

	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"index/internal/repository"
	"index/internal/database"

	"golang.org/x/crypto/bcrypt"
)

type User struct{
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

var masterKey = []byte(os.Getenv("KEY"))


func (u *User) SaveData(name, email, password string) error {
	validPassword, message := VerifyPassword(password)
	if !validPassword {
		return errors.New(message)
	}

	hash, err := HashPassword(password)
	if err != nil && validPassword {
		return err
	}

	return u.Repository.NewMysqlInsert(name, email, string(hash))
}


func (v *VerifyLogin) VerifyLoginFunction(name, email, password string) error {
	return v.Repository.NewVerifyLogin(name, email, password)

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
		_, err = database.Connect().Exec("INSERT INTO reset_password(user_id, token, expires_at) VALUES(?, ?, ?)", id, token, expiresAt)
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
	err = repository.Testing(hash, email)
	if err != nil {
		return err
	}

	StartToRemoverExpiredTokens()
	return nil
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


func Encrypt(text string) ([]byte, error) {
	block, _ := aes.NewCipher(masterKey)
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	return gcm.Seal(nonce, nonce, []byte(text), nil), nil
}


func Decrypt(cipherText []byte) (string, error) {
	block, _ := aes.NewCipher(masterKey)
	gcm, _ := cipher.NewGCM(block)
	nonceSize := gcm.NonceSize()
	nonce, actualCipher := cipherText[:nonceSize], cipherText[nonceSize:]
	text, err := gcm.Open(nil, nonce, actualCipher, nil)

	return string(text), err
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

	go func() {
		for range ticker.C {
			repository.RemoveExpiredToken(database.Connect())
		}
	}()

}
