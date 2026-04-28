package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	
	"ShieldAuth-API/internal/database"
	"ShieldAuth-API/internal/domain"

	"github.com/go-sql-driver/mysql"
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
}
type ValidTokenStruct struct {
	Database *sql.DB
}
type DeleteAccountStruct struct {
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
func NewValidTokenStruct(database *sql.DB) *ValidTokenStruct {
	return &ValidTokenStruct{
		Database: database,
	}
}
func NewDeleteAccountStruct(database *sql.DB) *DeleteAccountStruct {
	return &DeleteAccountStruct{
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


func(request *RequestStruct) GetID(ctx context.Context, u domain.User) (*domain.User, error) {
	user := &domain.User{}

	query := "SELECT id, email FROM users WHERE email = ? LIMIT 1"
	err := request.Database.QueryRowContext(ctx, query, u.Email).Scan(&user.Id, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return user, nil
}


func (resetPassword *ResetPasswordStruct) ResetPassword(ctx context.Context, currentPassword, newPassword, confirmNewPassword string) (error, string) {
	// var id int
	// var verify string
	// var token string
	// var used bool

	// tx, err := resetPassword.Database.BeginTx(ctx, nil)
	// if err != nil {
	// 	return err, ""
	// }
	// defer tx.Rollback()

	// if newPassword != confirmNewPassword {
	// 	return errors.New("ERROR"), ""
	// } else if currentPassword == newPassword {
	// 	return errors.New("ERROR"), ""
	// }

	// err = tx.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = ?", email).Scan(&id, &verify)
	// if err != nil {
	// 	return err, ""
	// }

	// passwordHash := bcrypt.CompareHashAndPassword([]byte(verify), []byte(currentPassword))
	// if passwordHash != nil {
	// 	return passwordHash, ""
	// }

	// err = tx.QueryRowContext(ctx, "SELECT token FROM reset_password WHERE user_id = ?", id).Scan(&token)
	// if err != nil {
	// 	return err, ""
	// }

	// err = tx.QueryRowContext(ctx, "SELECT used FROM reset_password WHERE token = ?", token).Scan(&used)
	// if err != nil {
	// 	return err, ""
	// }

	// allowed, err := LimitOfAttempts(ctx, email)
	// if err != nil {
	// 	return err, ""
	// }

	// if !allowed {
	// 	return errors.New("ERROR"), ""
	// }

	// if used {
	// 	return errors.New("ERROR"), ""
	// }

	// err = tx.Commit()
	// if err != nil {
	// 	return fmt.Errorf("Repository error: %w", err), ""
	// }

	// return nil, email
	return nil, ""
}


func (validToken *ValidTokenStruct) ValidToken(ctx context.Context, token, secondToken string) error {
	err := validToken.Database.QueryRowContext(ctx, "SELECT token FROM reset_password WHERE token = ? AND used = FALSE", token).Scan(&secondToken)
	if err != nil {
		return err
	}

	if token != secondToken {
		return errors.New("Invalid token")
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


func UpdatePassword(ctx context.Context, hash, email string) error {
	var id int
	tx, err := database.Connect().BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("Repository error: %w", err)
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "SELECT id FROM users WHERE email = ?", email).Scan(&id)
	if err != nil {
		return err
	}

	_, err = tx.ExecContext(ctx, "UPDATE users SET password = ? WHERE email = ?", hash, email)
	if err != nil {
		return fmt.Errorf("Repository error: %w", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE reset_password SET used = true WHERE user_id = ?", id)
	if err != nil {
		return fmt.Errorf("Repository error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Repository error: %w", err)
	}
	return nil
}


func LimitOfAttempts(ctx context.Context, email string) (bool, error) {
	var emailCount int

	tx, err := database.Connect().BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "DELETE FROM attempts WHERE attempted_at < NOW() - INTERVAL 24 HOUR")
	if err != nil {
		return false, err
	}

	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM attempts WHERE email = ? AND attempted_at > NOW() - INTERVAL 1 HOUR", email).Scan(&emailCount)
	if err != nil {
		return false, err
	}

	if emailCount >= 5 {
		return false, errors.New("ERROR")
	}

	_, err = tx.ExecContext(ctx, "INSERT INTO attempts (email, attempted_at) VALUES (?, NOW())", email)
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}


func InsertInto(ctx context.Context, userID int, token string, expiresAt time.Time) error {
	_, err := database.Connect().ExecContext(ctx, "INSERT INTO reset_password (user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("Repository error: %w", err)
	}
	return nil
}


func InsertIntoLoginAttempts(ctx context.Context, name, email string) error {

	_, err := database.Connect().ExecContext(ctx, "DELETE FROM login_attempts WHERE (email = ? OR name = ?) AND attempt_in < NOW() - INTERVAL 1 DAY", email, name)
	if err != nil {
		return err
	}

	_, err = database.Connect().ExecContext(ctx, "INSERT INTO login_attempts(name, email, attempt_in, success) VALUES(?, ?, NOW(), FALSE)", name, email)
	if err != nil {
		return err
	}
	return nil
}
