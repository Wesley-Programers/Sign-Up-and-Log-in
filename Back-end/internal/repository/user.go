package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
	
	"index/Back-end/internal/database"

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

var test = make([]string, 0, 1)

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
	_, err := register.Database.ExecContext(ctx, "INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, password)
	if err != nil {
		return fmt.Errorf("Repository error: %w", err)
	}
	return err
}


func (verifyLogin *VerifyLoginStruct) VerifyLogin(ctx context.Context, name, email, password string) error {

	var passwordHash string
	err := verifyLogin.Database.QueryRowContext(ctx, "SELECT password FROM users WHERE name = ? OR email = ?", name, email).Scan(&passwordHash)
	if err != nil {
		return errors.New("WRONG EMAIL OR NAME")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return errors.New("WRONG PASSWORD")
	}

	return nil
}


func (changeName *ChangeNameStruct) ChangeName(ctx context.Context, currentName, newName string) error {
	var exist bool

	tx, err := changeName.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE name = ?)", currentName).Scan(&exist)
	if err != nil {
		return err
	}

	if newName != "" && exist && currentName != newName {
		_, err = tx.ExecContext(ctx, "UPDATE users SET name = ? WHERE name = ?", newName, currentName)
		if err != nil {
			return fmt.Errorf("Repository error: %w", err)
		}

	} else if !exist {
		return errors.New("NAME DOES NOT EXIST")
	} else if currentName == newName {
		return errors.New("NEW NAME IS NOT DIFFERENT")
	} else {
		return errors.New("SOME ERROR")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Repository error: %w", err)
	}

	return nil
}


func (changeEmail *ChangeEmailStruct) ChangeEmail(ctx context.Context, currentEmail, newEmail, confirmNewEmail, password string) error {

	tx, err := changeEmail.Database.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exist bool
	var verifyPassword string

	err = tx.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", currentEmail).Scan(&exist)
	if err != nil {
		return err
	}

	err = tx.QueryRowContext(ctx, "SELECT password FROM users WHERE email = ?", currentEmail).Scan(&verifyPassword)
	if err != nil {
		return err
	}

	passwordHash := bcrypt.CompareHashAndPassword([]byte(verifyPassword), []byte(password))

	if exist && newEmail == confirmNewEmail && passwordHash == nil && newEmail != currentEmail {
		_, err = tx.ExecContext(ctx, "UPDATE users SET email = ? WHERE email = ?", newEmail, currentEmail)
		if err != nil {
			return fmt.Errorf("Repository error: %w", err)
		}

	} else if !exist {
		return errors.New("EMAIL DOES NOT EXIST")
	} else if newEmail != confirmNewEmail {
		return errors.New("ERROR")
	} else if passwordHash != nil {
		return errors.New("INCORRECT PASSWORD")
	} else if newEmail == currentEmail {
		return errors.New("ERROR")
	} else {
		return errors.New("SOME ERROR")
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Repository error: %w", err)
	}

	return nil
}


func (request *RequestStruct) Request(ctx context.Context, email string) (error, int) {
	var id int
	var verify string
	var attempts int

	tx, err := request.Database.BeginTx(ctx, nil)
	if err != nil {
		return err, 0
	}
	defer tx.Rollback()

	err = tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM attempts WHERE email = ? AND attempted_at > NOW() - INTERVAL 1 HOUR", email).Scan(&attempts)
	if err != nil {
		return err, 0
	}

	if attempts >= 5 {
		return errors.New("ERROR"), 0
	}

	err = tx.QueryRowContext(ctx, "SELECT id, email FROM users WHERE email = ?", email).Scan(&id, &verify)
	if err != nil {
		return err, 0
	}

	test = append(test, email)
	err = tx.Commit()
	if err != nil {
		return err, 0
	}

	return nil, id
}


func (resetPassword *ResetPasswordStruct) ResetPassword(ctx context.Context, currentPassword, newPassword, confirmNewPassword string) (error, string) {
	var id int
	var verify string
	var token string

	tx, err := resetPassword.Database.BeginTx(ctx, nil)
	if err != nil {
		return err, ""
	}
	defer tx.Rollback()

	email := test[0]

	err = tx.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = ?", email).Scan(&id, &verify)
	if err != nil {
		return err, ""
	}

	err = tx.QueryRowContext(ctx, "SELECT token FROM reset_password WHERE user_id = ?", id).Scan(&token)
	if err != nil {
		return err, ""
	}

	var used bool

	err = tx.QueryRowContext(ctx, "SELECT used FROM reset_password WHERE token = ?", token).Scan(&used)
	if err != nil {
		return err, ""
	}

	allowed, err := LimitOfAttempts(ctx, email)
	if err != nil {
		return err, ""
	}

	if !allowed {
		return errors.New("ERROR"), ""
	}

	passwordHash := bcrypt.CompareHashAndPassword([]byte(verify), []byte(currentPassword))
	if passwordHash == nil && !used {
		return nil, email
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("Repository error: %w", err), ""
	}
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

	err = tx.QueryRowContext(ctx, "SELECT password FROM users WHERE = ?", email).Scan(&verifyPassword)
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


func RemoveExpiredToken() error {

	context, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := database.Connect().ExecContext(context, "DELETE FROM reset_password WHERE expires_at < NOW() OR used = TRUE")
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
