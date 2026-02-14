package handlers

import (
	"fmt"
	"database/sql"
	"log"
	"net/http"
	"time"

	"index/internal/service"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Service *service.User
}

type LoginHandler struct {
	Service *service.VerifyLogin
}

type ChangeNameHandler struct {
	Service *service.ChangeName
}

type ChangeEmailHandler struct {
	Service *service.ChangeEmail
}

type RequestHandler struct {
	Service *service.Request
}

type ResetPasswordHandler struct {
	Service *service.ResetPassword
}

var teste = make([]string, 0, 1)
var store = sessions.NewCookieStore([]byte("KEY_SESSION"))
// var masterKey = []byte(os.Getenv("KEY"))

func (handler *Handler) NewSignUpHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type: ", "application/json")
	// w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")

	err := handler.Service.SaveData(name, email, password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("VALID DATA"))

	}

	log.Println("SUCCESS")

}


func (login *LoginHandler) NewHandlerLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type: ", "application/json")
	// w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.SetFlags(log.Lshortfile)

	nameOrEmail := r.FormValue("nameEmail")
	password := r.FormValue("passwordLog")

	err := login.Service.VerifyLoginFunction(nameOrEmail, nameOrEmail, password)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("VALID DATA"))

	} else {
		log.Println("ERROR: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	log.Println("SUCCESS")

}


func (changeName *ChangeNameHandler) ChangeNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type: ", "application/json")
	// w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	currentName := r.FormValue("currentName")
	newName := r.FormValue("newName")

	err := changeName.Service.ChangeNameFunction(currentName, newName)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("VALID DATA"))

	} else {
		log.Println("ERROR: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	log.Println("SUCCESS")
}


func (changeEmail *ChangeEmailHandler) ChangeEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.SetFlags(log.Lshortfile)

	currentEmail := r.FormValue("currentEmail")
	newEmail := r.FormValue("newEmail")
	confirmNewEmail := r.FormValue("confirmNewEmail")
	password := r.FormValue("currentPassword")

	err := changeEmail.Service.ChangeEmailFunction(currentEmail, newEmail, confirmNewEmail, password)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("VALID DATA"))

	} else {
		log.Println("ERROR: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return

	}

	log.Println("SUCCESS")
}


func (requestHandler *RequestHandler) RequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.SetFlags(log.Lshortfile)

	email := r.FormValue("email")

	err, link := requestHandler.Service.RequestFunction(email)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(link))

	} else {
		log.Println("ERROR: ", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("SUCCESS")
}


func (resetPasswordHandler *ResetPasswordHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type", "text/plain")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.SetFlags(log.Lshortfile)

	if r.Method == http.MethodPost {

		currentPassword := r.FormValue("currentPassword")
		newPassword := r.FormValue("newPassword")
		confirmPassword := r.FormValue("confirmPassword")

		err := resetPasswordHandler.Service.ResetPasswordFunction(currentPassword, newPassword, confirmPassword)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("VALID DATA"))

		} else {
			log.Println("ERROR: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println("SUCCESS")
	}
}


func HandlerAuthentication(database *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

	}
}


func LogoutHandler(database *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/plain")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		session, _ := store.Get(r, "tokenSession")
		session.Options.MaxAge = -1
		session.Save(r, w)

		http.SetCookie(w, &http.Cookie{
			Name: "user_data",
			Value: "",
			Path: "/",
			MaxAge: -1,
		})

	}

}


func HandlerDeleteAccount(database *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/plain")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.SetFlags(log.Lshortfile)

		if r.Method == http.MethodPost {

			err := r.ParseMultipartForm(10 << 20)
			if err != nil {
				http.Error(w, "ERROR: ", http.StatusBadRequest)
				return
			}

			emailConfirm := r.FormValue("emailConfirm")
			passwordConfirm := r.FormValue("passwordConfirm")

			var email bool
			var passwordHash string
			
			query := "DELETE FROM users WHERE password = ?"
			queryEmailSelect := "SELECT EXISTS(SELECT 1 FROM users WHERE email = ?)"
			queryPasswordSelect := "SELECT password FROM users WHERE email = ?"

			queryEmailErr := database.QueryRow(queryEmailSelect, emailConfirm).Scan(&email)
			if queryEmailErr != nil {
				log.Println("ERROR: ", queryEmailErr)
			}

			queryPasswordErr := database.QueryRow(queryPasswordSelect, emailConfirm).Scan(&passwordHash)
			if queryPasswordErr != nil {
				log.Println("ERROR: ", queryPasswordErr)

			} else if queryPasswordErr == sql.ErrNoRows {
				fmt.Println("SOMETHING")
				return
			}

			passwordHashCompare := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(passwordConfirm))

			if passwordHashCompare == nil && email {
				fmt.Println("CORRECT PASSWORD")

				w.WriteHeader(http.StatusOK)
				w.Write([]byte("EVERYTHING VALID"))
				_, err := database.Exec(query, passwordHashCompare)

				if err != nil {
					log.Fatal("SOMETHING BAD: ", err)
				}
				
			} else if passwordHashCompare != nil {
				fmt.Println("WRONG PASSWORD")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("INCORRECT PASSWORD"))

			} else if !email {
				fmt.Println("INCORRECT EMAIL")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("INCORRECT EMAIl"))

			} else {
				log.Println("SOME ERROR")

				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("SOME ERROR"))
			}

		} else {
			http.Error(w, "ERROR: ", http.StatusMethodNotAllowed)
		}

	}

}


// func Reset(database *sql.DB) http.HandlerFunc {

// 	return func (w http.ResponseWriter, r *http.Request) {

// 		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 		w.Header().Set("Access-Control-Allow-Credentials", "true")
// 		w.Header().Set("Content-Type", "text/plain")

// 		if r.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}

// 		log.SetFlags(log.Lshortfile)

// 		if r.Method == http.MethodPost {
			
// 			var id int
// 			var verify string
// 			var token string

// 			password := r.FormValue("currentPassword")
// 			newPassword := r.FormValue("newPassword")
// 			confirmPassword := r.FormValue("confirmPassword")

// 			email := teste[0]

// 			queryUpdatePassword, updatePasswordErro := database.Prepare("UPDATE users SET password = ? WHERE email = ?")
// 			if updatePasswordErro != nil {
// 				log.Println("ERROR: ", updatePasswordErro)
// 			}

// 			query := "SELECT id, password FROM users WHERE email = ?"
// 			queryPasswordError := database.QueryRow(query, email).Scan(&id, &verify)
// 			if queryPasswordError != nil {
// 				if queryPasswordError == sql.ErrNoRows {
// 					log.Println("PASSWORD NOT FOUND")
// 					return
// 				}
// 				log.Println("ERROR: ", queryPasswordError)
// 				return
// 			}

// 			getToken := "SELECT token FROM password_token WHERE user_id = ?"
// 			getTokenError := database.QueryRow(getToken, id).Scan(&token)
// 			if getTokenError != nil {
// 				if getTokenError == sql.ErrNoRows {
// 					log.Println("ID NOT FOUND")
// 					return
// 				}
// 				log.Println("ERROR: ", getTokenError)
// 				return
// 			}

// 			var idToken int
// 			var user_id int
// 			var expires_at string
// 			var used bool

// 			getInformation := "SELECT id, user_id, expires_at, used FROM password_token WHERE token = ?"
// 			getInformationError := database.QueryRow(getInformation, token).Scan(&idToken, &user_id, &expires_at, &used)
// 			if getInformationError != nil {
// 				if getInformationError == sql.ErrNoRows {
// 					log.Println("NOT FOUND")
// 					return
// 				}
// 				log.Println("ERROR: ", getInformationError)
// 				return
// 			}

// 			allowed, err := repository.LimitOfAttempts(database, email)
// 			if err != nil {
// 				log.Println("ERROR: ", err)
// 				http.Error(w, "INTERNAL ERROR", http.StatusInternalServerError)
// 				return
// 			}

// 			if !allowed {
// 				log.Println("A LOT OF ATTEMPTS, TRY AGAIN IN 15 MINUTES")
// 				http.Error(w, "LOT OF ATTEMPTS", http.StatusTooManyRequests)
// 				return
// 			}

// 			validPassword, message := service.VerifyPassword(newPassword)
// 			if !validPassword {
// 				if err := repository.AttemptLogs(database, email); err != nil {
// 					log.Println("ERROR: ", err)
// 				}

// 				w.Write([]byte(message))
// 				http.Error(w, "ERROR: ", http.StatusBadRequest)
// 				return

// 			}
// 			passwordHash := bcrypt.CompareHashAndPassword([]byte(verify), []byte(password))


// 			if passwordHash == nil && newPassword != "" && password != newPassword && utf8.RuneCountInString(newPassword) >= 8 && newPassword == confirmPassword && allowed && validPassword {
// 				fmt.Println("VALID PASSWORD")

// 				hash, err := service.HashPassword(newPassword)
// 				if err != nil {
// 					log.Fatal("ERROR: ", err)
// 					return
// 				}

// 				_, errorPassword := queryUpdatePassword.Exec(hash, teste[0])
// 				if errorPassword != nil {
// 					log.Fatal("ERROR: ", errorPassword)
// 					return
// 				}

// 				_, tokenUsedError := database.Exec("UPDATE password_token SET used = TRUE WHERE user_id = ?", user_id)
// 				if tokenUsedError != nil {
// 					log.Fatal("ERROR: ", tokenUsedError)
// 					return 
// 				}

// 				// repository.StartToRemoverExpiredTokens(database)

// 				ctx, cancel := context.WithCancel(context.Background())

// 				go func(ctx context.Context) {
// 					select {
// 					case <-time.After(20 * time.Second):
// 						teste = teste[:0]
// 						fmt.Println("SECOND EMPTY SLICE")

// 					case <-ctx.Done():
// 						fmt.Println("GOROTINE STOPPED")
// 					}
					
// 				}(ctx)
				
// 				w.WriteHeader(http.StatusOK)
// 				w.Write([]byte("VALID PASSWORD"))

// 				cancel()

// 			} else if passwordHash != nil {

// 				if err := repository.AttemptLogs(database, email); err != nil {
// 					log.Println("ERROR: ", err)
// 				}
// 				fmt.Println("WRONG PASSWORD")

// 				w.WriteHeader(http.StatusBadRequest)
// 				w.Write([]byte("INCORRECT PASSWORD"))

// 			} else if password == newPassword {

// 				if err := repository.AttemptLogs(database, email); err != nil {
// 					log.Println("ERROR: ", err)
// 				}
// 				fmt.Println("THE SAME PASSWORD")

// 				w.WriteHeader(http.StatusBadRequest)
// 				w.Write([]byte("THE PASSWORD ARE THE SAME"))

// 			} else if utf8.RuneCountInString(newPassword) < 8 {

// 				if err := repository.AttemptLogs(database, email); err != nil {
// 					log.Println("ERROR: ", err)
// 				}
// 				fmt.Println("SHORT PASSWORD")

// 				w.WriteHeader(http.StatusBadRequest)
// 				w.Write([]byte("SHORT PASSWORD"))

// 			} else if newPassword != confirmPassword {

// 				if err := repository.AttemptLogs(database, email); err != nil {
// 					log.Println("ERROR: ", err)
// 				}
// 				fmt.Println("PASSWORD CONFIRMATION IS WRONG")

// 				w.WriteHeader(http.StatusBadRequest)
// 				w.Write([]byte("PASSWORD CONFIRMATION IS WRONG"))

// 			} else {

// 				if err := repository.AttemptLogs(database, email); err != nil {
// 					log.Println("ERROR: ", err)
// 				}
// 				log.Println("SOME ERROR")

// 				w.WriteHeader(http.StatusBadRequest)
// 				w.Write([]byte("SOME ERROR"))
// 			}

// 		} else {
// 			http.Error(w, "METHOD NOT PERMITED", http.StatusMethodNotAllowed)
// 		}

// 	}
	
// }


func ValidTokenHandler(database *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "http://127.0.0.1:5500")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Content-Type", "text/plain")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		log.SetFlags(log.Lshortfile)

		if r.Method == http.MethodPost {

			var id int
			var expires_at string
			var used bool
	
			email := teste[0]
			err := database.QueryRow("SELECT id FROM users WHERE email = ?", email).Scan(&id)
			if err != nil {
				if err == sql.ErrNoRows {
					log.Println("EMAIL NOT FOUND")
					return
				}
				log.Println("ERROR: ", err)
				return
			}
	
			err = database.QueryRow("SELECT expires_at, used FROM password_token WHERE user_id = ?", id).Scan(&expires_at, &used)
			fmt.Println("EXPIRES AT: ", expires_at)

			layout := "2006-01-02 15:04:05"
			newExpiresAt, timeParseError := time.Parse(layout, expires_at)
			if timeParseError != nil {
				log.Println("ERROR: ", err)
				return
			}

			fmt.Println("NEW EXPIRES AT AGAIN: ", newExpiresAt)

			if err != nil {

				if err == sql.ErrNoRows || used == true || time.Now().After(newExpiresAt) {
					w.Write([]byte("INVALID TOKEN"))
					log.Println("INVALID TOKEN")
					return

				}
				log.Println("ERROR: ", err)
				return
	
			} else if err == nil {
				w.Write([]byte("VALID TOKEN"))
				time.Sleep(1000 * time.Millisecond)
				fmt.Println("USED: ", used)
				fmt.Println("VALID TOKEN")

			}

		}

	}

}
