package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"time"
	"unicode"

	"ShieldAuth-API/internal/domain"
	"ShieldAuth-API/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RequestReset(ctx context.Context, email string) (string, error)
	ValidToken(ctx context.Context, token string) (string, error)
}

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
type service struct {
	userRepo repository.UserRepository
	tokenRepo repository.ResetTokenRepository
}
type ResetPassword struct {
	Repository repository.ResetPassword
}
type DeleteAccount struct {
	Repository repository.DeleteAccount
}


func NewUserStruct(repository repository.User) *Register {
	return &Register{
		Repository: repository,
	}
}
func NewVerifyLogin(repository repository.LoginUser) *VerifyLogin {
	return &VerifyLogin{
		Repository: repository,
	}
}
func NewChangeName(repository repository.ChangeName) *ChangeName {
	return &ChangeName{
		Repository: repository,
	}
}
func NewChangeEmail(repository repository.ChangeEmail) *ChangeEmail {
	return &ChangeEmail{
		Repository: repository,
	}
}
func NewService(
	userRepo repository.UserRepository,
	tokenRepo repository.ResetTokenRepository,
) Service {
	return &service{
		userRepo: userRepo,
		tokenRepo: tokenRepo,
	}
}
func NewResetPassword(repository repository.ResetPassword) *ResetPassword {
	return &ResetPassword{
		Repository: repository,
	}
}
func NewDeleteAccount(repository repository.DeleteAccount) *DeleteAccount {
	return &DeleteAccount{
		Repository: repository,
	}
}


type ChangeEmailData struct {
	ID int
	CurrentEmail string
	NewEmail string
	ConfirmNewEmail string
	Password string
}
type ChangeNameData struct {
	ID int
	CurrentName string
	NewName string
	ConfirmNewName string
}
type LoginData struct {
	Name string
	Email string
	Password string
}
type RegisterData struct {
	Name string
	Email string
	Password string
}


func (register *Register) RegisterFunction(ctx context.Context, input RegisterData) error {
	validPassword, message := VerifyPassword(input.Password)
	if !validPassword {
		return errors.New(message)
	}

	hash, err := HashPassword(input.Password)
	if err != nil && validPassword {
		return err
	}

	err = register.Repository.Register(ctx, input.Name, input.Email, string(hash))
	if err != nil {
		return err
	}

	return nil
}


func (login *VerifyLogin) VerifyLoginFunction(ctx context.Context, input LoginData) (error, int) {

	identifier := input.Email
	if identifier == "" {
		identifier = input.Name
	}

	if identifier == "" {
		return domain.ErrInvalidData, 0
	}

	user, err := login.Repository.GetByCredentials(ctx, domain.User{Name: input.Name, Email: input.Email, PasswordHash: input.Password})
	if err != nil {
		return domain.ErrInvalidCredentials, 0
	}
	
	if err := user.PasswordValid(input.Password); err != nil {
		return domain.ErrInvalidPassword, 0
	}

	return nil, user.Id
}


func (changeName *ChangeName) ChangeNameFunction(ctx context.Context, input ChangeNameData) error {
	
	user, err := changeName.Repository.GetID(ctx, input.ID)
	if err != nil {
		return err
	}

	if err := user.ChangeName(input.CurrentName, input.NewName); err != nil {
		return err
	}

	return changeName.Repository.UpdateName(ctx, user)
}


func (changeEmail *ChangeEmail) ChangeEmailFunctionTest(ctx context.Context, input ChangeEmailData) error {

	user, err := changeEmail.Repository.GetID(ctx, input.ID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	if err := user.PasswordValid(input.Password); err != nil {
		return domain.ErrInvalidPassword
	}

	if err := user.ChangeEmail(input.CurrentEmail, input.NewEmail, input.ConfirmNewEmail); err != nil {
		return domain.ErrInternal
	}

	return changeEmail.Repository.UpdateEmail(ctx, user)
}


func (s *service) RequestReset(ctx context.Context, email string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil
	}

	token, err := GenerateTokens()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	hash := tokenHash(token)
	expiresAt := time.Now().Add(15 * time.Minute)

	if err := s.tokenRepo.InvalidateAll(ctx, user.Id); err != nil {
		return "", fmt.Errorf("invalidate token: %w", err)
	}

	if err := s.tokenRepo.Save(ctx, user.Id, hash, expiresAt); err != nil {
		return "", fmt.Errorf("save token: %w", err)
	}

	return token, nil
}


func (s *service) ValidToken(ctx context.Context, token string) (string, error) {
	if token == "" {
		return "", domain.ErrInvalidToken
	}

	hash := tokenHash(token)
	userID, err := s.tokenRepo.FindValid(ctx, hash)
	if err != nil {
		return "", domain.ErrInvalidToken
	}

	if err := s.tokenRepo.MarkUsed(ctx, hash); err != nil {
		return "", fmt.Errorf("mark token used: %w", err)
	}
	
	return userID, nil
}


func (resetPasword *ResetPassword) ResetPasswordFunction(ctx context.Context, currentPassword, newPassword, confirmNewPassword string) error {
	validPassword, message := VerifyPassword(newPassword)
	if !validPassword {
		return errors.New(message)
	}

	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}

	err, email := resetPasword.Repository.ResetPassword(ctx, currentPassword, newPassword, confirmNewPassword)
	if err != nil {
		_, err = repository.LimitOfAttempts(ctx, email)
		if err != nil {
			return err
		}
		return err
	}
	
	err = repository.UpdatePassword(ctx, hash, email)
	if err != nil {
		_, err = repository.LimitOfAttempts(ctx, email)
		if err != nil {
			return err
		}
		return err
	}

	return nil
}


func (delete *DeleteAccount) DeleteAccountFunction(ctx context.Context, email, password string) error {
	return delete.Repository.DeleteAccount(ctx, email, password)
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
	return base64.URLEncoding.EncodeToString(bytes), nil
	// return hex.EncodeToString(bytes), nil
}


func StartToRemoverExpiredTokens(database *sql.DB) {
	ticker := time.NewTicker(20 * time.Second)

	for range ticker.C {
		repository.RemoveExpiredToken(database)
	}
}


func tokenHash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}


func TokenJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_KEY")
	return tokenJwt.SignedString([]byte(secretKey))
}
