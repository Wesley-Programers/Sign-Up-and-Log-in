package handlers

import (
	"fmt"
	"database/sql"
	"log"
	"net/http"
	"time"

	"index/internal/service"

	"github.com/gorilla/sessions"
	// "golang.org/x/crypto/bcrypt"
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

type DeleteAccountHandler struct {
	Service *service.DeleteAccount
}

var teste = make([]string, 0, 1)
var store = sessions.NewCookieStore([]byte("KEY_SESSION"))
// var masterKey = []byte(os.Getenv("KEY"))

func (handler *Handler) NewSignUpHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.Header().Set("Content-Type: ", "application/json")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {

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

	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
	}

}


func (login *LoginHandler) NewHandlerLogin(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.SetFlags(log.Lshortfile)

	if r.Method == http.MethodPost {
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

	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
	}

}


func (changeName *ChangeNameHandler) ChangeNameHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
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

	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
	}

}


func (changeEmail *ChangeEmailHandler) ChangeEmailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.SetFlags(log.Lshortfile)

	if r.Method == http.MethodPost {
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

	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
	}

}


func (requestHandler *RequestHandler) RequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.SetFlags(log.Lshortfile)

	if r.Method == http.MethodPost {
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

	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
	}

}


func (resetPasswordHandler *ResetPasswordHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Content-Type", "text/plain")

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

	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
	}

}


func (deleteAccountHandler *DeleteAccountHandler) DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.SetFlags(log.Lshortfile)

	if r.Method == http.MethodPost {

		email := r.FormValue("emailConfirm")
		password := r.FormValue("passwordConfirm")

		err := deleteAccountHandler.Service.DeleteAccountFunction(email, password)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("VALID"))

		} else {
			log.Println("ERROR")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println("SUCCESS")

	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
	}
}


func ValidTokenHandler(database *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {


		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

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
	
			} else {
				w.Write([]byte("VALID TOKEN"))
				time.Sleep(1000 * time.Millisecond)
				fmt.Println("USED: ", used)
				fmt.Println("VALID TOKEN")

			}
		}
	}
}
