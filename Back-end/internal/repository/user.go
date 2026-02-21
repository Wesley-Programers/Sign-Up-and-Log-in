package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"index/internal/database"

	"golang.org/x/crypto/bcrypt"
)


type RegisterStruct struct{
	Database *sql.DB
}
type VerifyLoginStruct struct{}
type ChangeNameStruct struct{}
type ChangeEmailStruct struct{}
type RequestStruct struct{}
type ResetPasswordStruct struct{}
type ValidTokenStruct struct{}
type DeleteAccountStruct struct{}

var test = make([]string, 0, 1)

func NewRegisterStruct(database *sql.DB) *RegisterStruct {
	return &RegisterStruct{
		Database: database,
	}
}


func (register *RegisterStruct) Register(ctx context.Context, name, email, password string) error {
	_, err := register.Database.ExecContext(ctx, "INSERT INTO users (name, email, password) VALUES (?, ?, ?)", name, email, password)
	if err != nil {
		return fmt.Errorf("REPOSITORY ERROR: %w", err)
	}
	return err
}


func (verifyLogin *VerifyLoginStruct) VerifyLogin(name, email, password string) error {

	var passwordHash string
	err := database.Connect().QueryRow("SELECT password FROM users WHERE name = ? OR email = ?", name, email).Scan(&passwordHash)
	if err != nil {
		return errors.New("WRONG EMAIL OR NAME")
	}

	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	if err != nil {
		return errors.New("WRONG PASSWORD")
	}
	
	return nil
}


func (changeName *ChangeNameStruct) ChangeName(currentName, newName string) error {
	var exist bool

	tx, err := database.Connect().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE name = ?)", currentName).Scan(&exist)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	if newName != "" && exist && currentName != newName {
		_, err = tx.Exec("UPDATE users SET name = ? WHERE name = ?", newName, currentName)
		if err != nil {
			log.Println("ERROR: ", err)
			return err
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
		return err
	}

	return nil
}


func (changeEmail *ChangeEmailStruct) ChangeEmail(currentEmail, newEmail, confirmNewEmail, password string) error {
	
	tx, err := database.Connect().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	
	var exist bool
	var verifyPassword string

	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", currentEmail).Scan(&exist)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	err = tx.QueryRow("SELECT password FROM users WHERE email = ?", currentEmail).Scan(&verifyPassword)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	passwordHash := bcrypt.CompareHashAndPassword([]byte(verifyPassword), []byte(password))

	if exist && newEmail == confirmNewEmail && passwordHash == nil && newEmail != currentEmail {
		_, err = tx.Exec("UPDATE users SET email = ? WHERE email = ?", newEmail, currentEmail)
		if err != nil {
			log.Println("ERROR: ", err)
			return err
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
		return err
	}

	return nil
}


func (request *RequestStruct) Request(email string) (error, int) {
	var id int
	var verify string

	err := database.Connect().QueryRow("SELECT id, email FROM users WHERE email = ?", email).Scan(&id, &verify)
	if err != nil {
		log.Println("ERROR: ", err)
		return err, 0
	}

	if verify == email {
		test = append(test, email)
		return nil, id
	}

	return nil, 0
}


func (resetPassword *ResetPasswordStruct) ResetPassword(currentPassword, newPassword, confirmNewPassword string) (error, string) {
	var id int
	var verify string
	var token string

	tx, err := database.Connect().Begin()
	if err != nil {
		return err, ""
	}
	defer tx.Rollback()

	email := test[0]

	err = tx.QueryRow("SELECT id, password FROM users WHERE email = ?", email).Scan(&id, &verify)
	if err != nil {
		log.Println("ERROR: ", err)
		return err, ""
	}

	err = tx.QueryRow("SELECT token FROM reset_password WHERE user_id = ?", id).Scan(&token)
	if err != nil {
		log.Println("ERROR: ", err)
		return err, ""
	}

	var used bool

	err = tx.QueryRow("SELECT used FROM reset_password WHERE token = ?", token).Scan(&used)
	if err != nil {
		log.Println("ERROR: ", err)
		return err, ""
	}

	allowed, err := LimitOfAttempts(email)
	if err != nil {
		log.Println("ERROR: ", err)
		return err, ""
	}

	if !allowed {
		return errors.New("ERROR"), ""
	}

	passwordHash := bcrypt.CompareHashAndPassword([]byte(verify), []byte(currentPassword))
	if passwordHash == nil && !used {
		
		ctx, cancel := context.WithCancel(context.Background())
		go func(ctx context.Context) {
			select {
			case <-time.After(15 *time.Second):
				test = test[:0]

			case <-ctx.Done():
				log.Println("")
				return
			}
		}(ctx)
		
		cancel()
		return nil, email
	}

	err = tx.Commit()
	if err != nil {
		return err, ""
	}
	return nil, ""
}


func Test(token, secondToken string) error {
	err := database.Connect().QueryRow("SELECT token FROM reset_password WHERE token = ? AND used = FALSE AND expires_at > NOW()", token).Scan(&secondToken)
	if err != nil {
		return err
	}
	return nil
}


func (deleteAccount *DeleteAccountStruct) DeleteAccount(email, password string) error {
	
	tx, err := database.Connect().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var verifyEmail bool
	var verifyPassword string

	err = tx.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)", email).Scan(&verifyEmail)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	err = tx.QueryRow("SELECT password FROM users WHERE = ?", email).Scan(&verifyPassword)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	passwordHash := bcrypt.CompareHashAndPassword([]byte(verifyPassword), []byte(password))

	if passwordHash == nil && verifyEmail {
		_, err = tx.Exec("DELETE FROM users WHERE password = ?", passwordHash)
		if err != nil {
			log.Println("ERROR: ", err)
			return err
		}

	} else {
		log.Println("ERROR: ", passwordHash)
		return passwordHash
	}

	err = tx.Commit()
	if err != nil {
		log.Println("ERROR: ", err)
		return err
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


func UpdatePassword(hash, email string) error {
	var id int
	tx, err := database.Connect().Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	err = tx.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&id)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE users SET password = ? WHERE email = ?", hash, email)
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE reset_password SET used = true WHERE user_id = ?", id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}


func LimitOfAttempts(email string) (bool, error) {
	var emailCount int

	tx, err := database.Connect().Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM attempts WHERE attempted_at < NOW() - INTERVAL 24 HOUR")
	if err != nil {
		return false, err
	}

	err = tx.QueryRow("SELECT COUNT(*) FROM attempts WHERE email = ? AND attempted_at > NOW() - INTERVAL 1 HOUR", email).Scan(&emailCount)
	if err != nil {
		return false, err
	}

	if emailCount >= 5 {
		return false, errors.New("ERROR")
	}

	_, err = tx.Exec("INSERT INTO attempts (email, attempted_at) VALUES (?, NOW())", email)
	if err != nil {
		return false, err
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}


func InsertInto(userID int, token string, expiresAt time.Time) error {
	_, err := database.Connect().Exec("INSERT INTO reset_password (user_id, token, expires_at) VALUES (?, ?, ?)", userID, token, expiresAt)
	if err != nil {
		return err
	}
	return nil
}
