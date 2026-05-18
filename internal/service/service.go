package service

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"ShieldAuth-API/internal/domain"
	"ShieldAuth-API/internal/repository"
	"ShieldAuth-API/internal/security"
)

type Service interface {
	RequestReset(ctx context.Context, email string) (string, error)
	ValidToken(ctx context.Context, token string) error
}
type Security interface {
	GenerateToken() (string, error)
	TokenHash(token string) string
}
type Limiter interface {
	CheckLimit(ctx context.Context, key string, maxAttempts int, window time.Duration) (bool, error)
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
type ResetStore interface {
	Get(ctx context.Context, token string) (string, error)
	Delete(ctx context.Context, token string) error
	Save(ctx context.Context, token string, userID int, ttl time.Duration) error
}

type service struct {
	userRepo repository.UserRepository
	Security ResetStore
	security Security
	Limiter  Limiter
}
type ResetPassword struct {
	Repository repository.ResetPassword
	Security   *security.ResetPassword
}
type DeleteAccount struct {
	Repository repository.DeleteAccount
}
type ValidTesting struct {
	Repository repository.ResetTokenRepository
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
	security Security,
	resetStore ResetStore,
	limiter Limiter,
) Service {
	return &service{
		userRepo: userRepo,
		security: security,
		Security: resetStore,
		Limiter:  limiter,
	}
}
func NewValidToken(repository repository.ResetTokenRepository) *ValidTesting {
	return &ValidTesting{
		Repository: repository,
	}
}
func NewResetPassword(repository repository.ResetPassword, sec *security.ResetPassword) *ResetPassword {
	return &ResetPassword{
		Repository: repository,
		Security:   sec,
	}
}
func NewDeleteAccount(repository repository.DeleteAccount) *DeleteAccount {
	return &DeleteAccount{
		Repository: repository,
	}
}

type ChangeEmailData struct {
	ID             	int
	CurrentEmail    string
	NewEmail        string
	ConfirmNewEmail string
	Password        string
}
type ChangeNameData struct {
	ID             int
	CurrentName    string
	NewName        string
	ConfirmNewName string
}
type LoginData struct {
	Name     string
	Email    string
	Password string
}
type RegisterData struct {
	Name     string
	Email    string
	Password string
}
type ResetPasswordData struct {
	Token           string
	NewPasword      string
	ConfirmPassword string
}

func (register *Register) RegisterFunction(ctx context.Context, input RegisterData) error {
	validPassword, message := security.VerifyPassword(input.Password)
	if !validPassword {
		return errors.New(message)
	}

	hash, err := security.HashPassword(input.Password)
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

	normalizedEmail := strings.ToLower(strings.TrimSpace(email))
	sum := sha256.Sum256([]byte(normalizedEmail))
	key := fmt.Sprintf("forgot-password:email:%x", sum)

	allowed, err := s.Limiter.CheckLimit(ctx, fmt.Sprintf("reset-password: %s", key), 10, time.Minute)
	if err != nil {
		return "", fmt.Errorf("rate limit failed: %w", err)
	}

	if !allowed {
		return "", errors.New("too many attempts")
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", nil
	}

	token, err := s.security.GenerateToken()
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	s.Security.Save(ctx, token, user.Id, 15*time.Minute)

	return token, nil
}

func (s *service) ValidToken(ctx context.Context, token string) error {
	if token == "" {
		return domain.ErrInvalidToken
	}

	_, err := s.Security.Get(ctx, token)
	if err != nil {
		return domain.ErrInvalidToken
	}

	return nil
}

func (r *ResetPassword) ResetPasswordFunction(ctx context.Context, token string, input ResetPasswordData) error {
	if input.NewPasword != input.ConfirmPassword {
		return errors.New("password do not match")
	}

	valid, msg := security.VerifyPassword(input.NewPasword)
	if !valid {
		return errors.New(msg)
	}

	userID, err := r.Security.Get(ctx, token)
	if err != nil {
		return fmt.Errorf("invalid or expired token: %w", err)
	}

	defer r.Security.Delete(ctx, token)

	hashedPassword, err := security.HashPassword(input.NewPasword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = r.Repository.UpdatePassword(ctx, userID, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (delete *DeleteAccount) DeleteAccountFunction(ctx context.Context, email, password string) error {
	return delete.Repository.DeleteAccount(ctx, email, password)
}

func StartToRemoverExpiredTokens(database *sql.DB) {
	ticker := time.NewTicker(20 * time.Second)

	for range ticker.C {
		repository.RemoveExpiredToken(database)
	}
}
