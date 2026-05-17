package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"ShieldAuth-API/internal/domain"

	"github.com/go-sql-driver/mysql"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)


type RegisterStruct struct {
	Database *sql.DB
}
type VerifyLoginStruct struct {
	Database *sql.DB
}
type ChangeNameStruct struct {
	Database *sql.DB
}
type ChangeEmailStruct struct {
	Database *sql.DB
}
type RequestStruct struct {
	Database *sql.DB
}
type ResetPasswordStruct struct {
	Database *sql.DB
	redis *redis.Client
}
type ValidTokenStruct struct {
	Database *sql.DB
}
type DeleteAccountStruct struct {
	Database *sql.DB
}
type SessionAndAudit struct {
	Database *sql.DB
}


func NewRegisterStruct(database *sql.DB) *RegisterStruct {
	return &RegisterStruct{
		Database: database,
	}
}
func NewVerifyLoginStruct(database *sql.DB) *VerifyLoginStruct {
	return &VerifyLoginStruct{
		Database: database,
	}
}
func NewChangeNameStruct(database *sql.DB) *ChangeNameStruct {
	return &ChangeNameStruct{
		Database: database,
	}
}
func NewChangeEmailStruct(database *sql.DB) *ChangeEmailStruct {
	return &ChangeEmailStruct{
		Database: database,
	}
}
func NewRequestStruct(database *sql.DB) *RequestStruct {
	return &RequestStruct{
		Database: database,
	}
}
func NewResetPasswordStruct(database *sql.DB) *ResetPasswordStruct {
	return &ResetPasswordStruct{
		Database: database,
	}
}
func NewDeleteAccountStruct(database *sql.DB) *DeleteAccountStruct {
	return &DeleteAccountStruct{
		Database: database,
	}
}
func NewValidTokenStruct(database *sql.DB) *ValidTokenStruct {
	return &ValidTokenStruct{
		Database: database,
	}
}


func (register *RegisterStruct) Register(ctx context.Context, name, email, password string) error {
	u := &domain.User{
		Name: name,
		Email: email,
		PasswordHash: password,
	}

	_, err := register.Database.ExecContext(ctx, "INSERT INTO users (name, email, password) VALUES (?, ?, ?)", u.Name, u.Email, u.PasswordHash)
	if err != nil {
		if mySQLError, ok := errors.AsType[*mysql.MySQLError](err); ok {
			if mySQLError.Number == 1062 {
				return domain.ErrEmailAlreadyExists
			}
		}
		return fmt.Errorf("Repository error: failed to insert user: %w", err)
	}
	return nil
}


func (loginStruct *VerifyLoginStruct) GetByCredentials(ctx context.Context, u domain.User) (*domain.User, error) {
	user := &domain.User{}

	query := "SELECT id, name, email, password FROM users WHERE name = ? OR email = ? LIMIT 1"
	err := loginStruct.Database.QueryRowContext(ctx, query, u.Name, u.Email).Scan(&user.Id, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil {	
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}


func (changeName *ChangeNameStruct) GetID(ctx context.Context, id int) (*domain.User, error) {
	var test struct {
		ID int
		Name string
	}

	query := "SELECT id, name FROM users WHERE id = ?"
	err := changeName.Database.QueryRowContext(ctx, query, id).Scan(&test.ID, &test.Name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return domain.RestoreName(test.ID, test.Name), nil
}

func (changeName *ChangeNameStruct) UpdateName(ctx context.Context, user *domain.User) error {
	query := "UPDATE users SET name = ? WHERE id = ?"
	_, err := changeName.Database.ExecContext(ctx, query, user.NAME(), user.ID())
	return err
}


func (changeEmail *ChangeEmailStruct) GetID(ctx context.Context, id int) (*domain.User, error) {
	var test struct {
		ID int
		Email string
		PasswordHash string
	}

	query := "SELECT id, email, password FROM users WHERE id = ?"
	err := changeEmail.Database.QueryRowContext(ctx, query, id).Scan(&test.ID, &test.Email, &test.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	
	return domain.RestoreEmail(test.ID, test.Email, test.PasswordHash), nil
}


func (changeEmail *ChangeEmailStruct) UpdateEmail(ctx context.Context, user *domain.User) error {
	query := "UPDATE users SET email = ? WHERE id = ?"
	_, err := changeEmail.Database.ExecContext(ctx, query, user.EMAIL(), user.ID())
	return err
}


func(r *RequestStruct) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user := &domain.User{}

	err := r.Database.QueryRowContext(ctx, `SELECT id, email FROM users WHERE email = ?`, email).Scan(&user.Id, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return user, nil
}


func (v *ValidTokenStruct) Save(ctx context.Context, userID int, tokenHash string, expiresAt time.Time) error {
	_, err := v.Database.ExecContext(ctx, `INSERT INTO reset_password (user_id, token_hash, expires_at, used) VALUES (?, ?, ?, FALSE)`, userID, tokenHash, expiresAt)
	if err != nil {
		return fmt.Errorf("save reset token: %w", err)
	}

	return nil
}


func (v *ValidTokenStruct) FindValid(ctx context.Context, tokenHash string) (string, error) {
	var userID string

	err := v.Database.QueryRowContext(ctx, `SELECT user_id FROM reset_password WHERE token_hash = ? AND used = FALSE AND expires_at > ?`, tokenHash, time.Now()).Scan(&userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", domain.ErrInvalidToken
		}
		return "", fmt.Errorf("find valid token: %w", err)
	}

	return userID, nil
}


func (r *ResetPasswordStruct) AllowReset(ctx context.Context, email string) error {

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	key := fmt.Sprintf("reset_attempts:%s", email)

	count, err := r.redis.Incr(ctx, key).Result()
	if err != nil {
		return fmt.Errorf("redis incr failed: %w", err)
	}

	if count == 1 {
		err = r.redis.Expire(ctx, key, time.Hour).Err()
		if err != nil {
			return fmt.Errorf("redis expire failed: %w", err)
		}
	}

	if count > 5 {
		return fmt.Errorf("too many attempts")
	}

	return nil
}


func (deleteAccount *DeleteAccountStruct) DeleteAccount(ctx context.Context, email, password string) error {

	tx, err := deleteAccount.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var verifyEmail bool
	var verifyPassword string

	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&verifyEmail)
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(ctx, "SELECT password FROM users WHERE email = ?", email).Scan(&verifyPassword)
	if err != nil {
		return err
	}

	passwordHash := bcrypt.CompareHashAndPassword([]byte(verifyPassword), []byte(password))

	if passwordHash == nil && verifyEmail {
		_, err = tx.ExecContext(ctx, "DELETE FROM users WHERE password = ?", passwordHash)
		if err != nil {
			return fmt.Errorf("Repository error: %w", err)
		}

	} else {
		log.Println("ERROR: ", passwordHash)
		return passwordHash
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Repository error: %w", err)
	}
	return nil
}


func RemoveExpiredToken(database *sql.DB) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := database.ExecContext(ctx, "DELETE FROM reset_password WHERE expires_at < NOW() OR used = TRUE")
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows > 0 {
		log.Printf("REMOVED %d expired or used tokens", rows)
	}

	return nil
}


func (r *ResetPasswordStruct) UpdatePassword(ctx context.Context, userID string, passwordHash string) error {
	
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.Database.ExecContext(ctx, `UPDATE users SET password_hash = ? WHERE id = ?`, passwordHash, userID)
	if err != nil {
		return fmt.Errorf("update password failed: %w", err)
	}

	return nil
}


func (s *SessionAndAudit) InsertIntoLoginAudits(ctx context.Context, email string, success bool, failureReason string) error {
	
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	_, err := s.Database.ExecContext(ctx, `INSERT INTO login_attempts_audit (email, success, failure_reason, attempted_at) VALUES (?, ?, ?, NOW())`, email, success, failureReason)
	if err != nil {
		return fmt.Errorf("insert login audit failed: %w", err)
	}

	return nil
}


func (s *SessionAndAudit) CreateSession(ctx context.Context, userID int, refreshTokenHash string) error {
	_, err := s.Database.ExecContext(ctx, `INSERT INTO sessions (user_id, refresh_token_hash, revoked, expires_at, created_at) VALUES (?, ?, false, DATE_ADD(NOW(), INTERNAL 7 DAY), NOW())`, userID, refreshTokenHash)
	return err
}