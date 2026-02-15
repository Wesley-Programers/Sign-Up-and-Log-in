package repository

import (
	"errors"
	"log"

	"index/internal/database"

	"golang.org/x/crypto/bcrypt"
)

type Register struct{}
type VerifyLoginStruct struct{}
type ChangeNameStruct struct{}
type ChangeEmailStruct struct{}
type RequestStruct struct{}
type ResetPasswordStruct struct{}
type DeleteAccountStruct struct{}

var test = make([]string, 0, 1)


func (register *Register) NewMysqlInsert(name, email, password string) error {
	query := "INSERT INTO users (name, email, password) VALUES (?, ?, ?);"

	_, erroInsert := database.Connect().Exec(query, name, email, password)
	return erroInsert
}


func (verifyLogin *VerifyLoginStruct) NewVerifyLogin(name, email, password string) error {
	query := "SELECT password FROM users WHERE name = ? OR email = ?"

	var passwordHash string
	err := database.Connect().QueryRow(query, name, email).Scan(&passwordHash)
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

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE name = ?)"
	queryUpdateName, err := database.Connect().Prepare("UPDATE users SET name = ? WHERE name = ?")

	err = database.Connect().QueryRow(query, currentName).Scan(&exist)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	if newName != "" && exist && currentName != newName {
		_, err := queryUpdateName.Exec(newName, currentName)
		if err != nil {
			log.Println("ERROR: ", err)
			return err
		}

		return nil

	} else if !exist {
		return errors.New("NAME DOES NOT EXIST")
	} else if currentName == newName {
		return errors.New("NEW NAME IS NOT DIFFERENT")
	} else {
		return errors.New("SOME ERROR")
	}
}


func (changeEmail *ChangeEmailStruct) ChangeEmail(currentEmail, newEmail, confirmNewEmail, password string) error {
	var exist bool
	var verifyPassword string

	query := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"
	queryPassword := "SELECT password FROM users WHERE email = ?"
	queryUpdateEmail, err := database.Connect().Prepare("UPDATE users SET email = ? WHERE email = ?")
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	err = database.Connect().QueryRow(query, currentEmail).Scan(&exist)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	err = database.Connect().QueryRow(queryPassword, currentEmail).Scan(&verifyPassword)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	passwordHash := bcrypt.CompareHashAndPassword([]byte(verifyPassword), []byte(password))
	if exist && newEmail == confirmNewEmail && passwordHash == nil && newEmail != currentEmail {
		_, err = queryUpdateEmail.Exec(newEmail, currentEmail)
		if err != nil {
			log.Println("ERROR: ", err)
			return err
		}
		return nil

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
}


func (request *RequestStruct) Request(email string) (error, int) {
	var id int
	var verify string

	query := "SELECT id, email FROM users WHERE email = ?"
	err := database.Connect().QueryRow(query, email).Scan(&id, &verify)
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

	email := test[0]
	var err error

	querySelect := "SELECT id, password FROM users WHERE email = ?"
	err = database.Connect().QueryRow(querySelect, email).Scan(&id, &verify)
	if err != nil {
		log.Println("ERROR: ", err)
		return err, ""
	}

	queryToken := "SELECT token FROM reset_password WHERE user_id = ?"
	err = database.Connect().QueryRow(queryToken, id).Scan(&token)
	if err != nil {
		log.Println("ERROR: ", err)
		return err, ""
	}

	var idToken int
	var user_id int
	var expires_at string
	var used bool

	queryReset := "SELECT id, user_id, expires_at, used FROM reset_password WHERE token = ?"
	err = database.Connect().QueryRow(queryReset, token).Scan(&idToken, user_id, expires_at, used)
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
		// ctx, cancel := context.WithCancel(context.Background())
		// go func(ctx context.Context) {
		// 	select {
		// 	case <-time.After(20 *time.Second):
		// 		test = test[:0]

		// 	case <-ctx.Done():
		// 		log.Println("")
		// 		return
		// 	}

		// }(ctx)
		
		// cancel()
		return nil, email
	}

	return nil, ""
}


func (deleteAccount *DeleteAccountStruct) DeleteAccount(email, password string) error {
	var verifyEmail bool
	var verifyPassword string

	query := "DELETE FROM users WHERE password = ?"
	queryEmail := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"
	queryPassword := "SELECT password FROM users WHERE email = ?"

	err := database.Connect().QueryRow(queryEmail, email).Scan(&verifyEmail)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	err = database.Connect().QueryRow(queryPassword, email).Scan(&verifyPassword)
	if err != nil {
		log.Println("ERROR: ", err)
		return err
	}

	passwordHash := bcrypt.CompareHashAndPassword([]byte(verifyPassword), []byte(password))
	if passwordHash == nil && verifyEmail {
		_, err = database.Connect().Exec(query, passwordHash)
		if err != nil {
			log.Println("ERROR: ", err)
			return err
		}

		return nil

	} else {
		log.Println("ERROR: ", passwordHash)
		return passwordHash
	}

}


func RemoveExpiredToken() {

	_, queryError := database.Connect().Exec("DELETE FROM password_token WHERE expires_at < NOW()")
	if queryError != nil {
		log.Println("ERROR: ", queryError)
		return
	}

	log.Println("TOKENS REMOVED")

}


func Testing(hash, email string) (error, int) {
	var id int

	query, err := database.Connect().Prepare("UPDATE users SET password = ? WHERE email = ?")
	if err != nil {
		return err, 0
	}

	_, err = query.Exec(hash, email)
	if err != nil {
		return err, 0
	}

	err = database.Connect().QueryRow("SELECT id WHERE email = ?", email).Scan(&id)
	if err != nil {
		return err, 0
	}

	return nil, id
}


func TokenUsed(id int) error {
	query, err := database.Connect().Prepare("UPDATE reset_password SET used = true WHERE user_id = ?")
	if err != nil {
		return err
	}

	_, err = query.Exec()
	if err != nil {
		return err
	}

	return nil
}


func LimitOfAttempts(email string) (bool, error) {

	var emailCount int

	err := database.Connect().QueryRow("SELECT COUNT(*) FROM attempts WHERE email = ?", email).Scan(&emailCount)
	if err != nil {
		return false, err
	}

	if emailCount >= 10 {
		return false, errors.New("ERROR")
	}

	return true, nil

}


func AttemptLogs(email string) error {
	_, err := database.Connect().Exec("INSERT INTO attempts(email) VALUES(?)", email)
	return err
}
