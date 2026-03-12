package handlers

import (
	"html/template"
	"errors"
	"log"
	"net/http"

	"index/Back-end/internal/domain"
	"index/Back-end/internal/service"
	"index/Back-end/internal/ui"
)

type RegisterHandler struct {
	Service *service.Register
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

type ValidTokenHandler struct {
	Service *service.ValidToken
}

func NewRegisterHanlder(service *service.Register) *RegisterHandler {
	return &RegisterHandler{
		Service: service,
	}
}
func NewLoginHandler(service *service.VerifyLogin) *LoginHandler {
	return &LoginHandler{
		Service: service,
	}
}
func NewChangeNameHandler(service *service.ChangeName) *ChangeNameHandler {
	return &ChangeNameHandler{
		Service: service,
	}
}
func NewChangeEmailHandler(service *service.ChangeEmail) *ChangeEmailHandler {
	return &ChangeEmailHandler{
		Service: service,
	}
}
func NewRequestHandler(service *service.Request) *RequestHandler {
	return &RequestHandler{
		Service: service,
	}
}
func NewResetPasswordHandler(service *service.ResetPassword) *ResetPasswordHandler {
	return &ResetPasswordHandler{
		Service: service,
	}
}
func NewDeleteAccountHandler(service *service.DeleteAccount) *DeleteAccountHandler {
	return &DeleteAccountHandler{
		Service: service,
	}
}
func NewValidTokenHandler(service *service.ValidToken) *ValidTokenHandler {
	return &ValidTokenHandler{
		Service: service,
	}
}

var tmpl = template.Must(template.ParseFS(ui.Files, "templates/reset.html"))

func (handler *RegisterHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {

		ctx := r.Context()

		name := r.FormValue("name")
		email := r.FormValue("email")
		password := r.FormValue("password")
	
		err := handler.Service.RegisterFunction(ctx, name, email, password)
		if err != nil {
			switch {
			case errors.Is(err, domain.ErrEmailAlreadyExist):
				http.Error(w, err.Error(), http.StatusConflict)

			case errors.Is(err, domain.ErrInvaliData):
				http.Error(w, err.Error(), http.StatusBadRequest)

			default:
				log.Printf("Unexpected error: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusCreated)
	
	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
	}
}


func (login *LoginHandler) HandlerLogin(w http.ResponseWriter, r *http.Request) {

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

		ctx := r.Context()
	
		err := login.Service.VerifyLoginFunction(ctx, nameOrEmail, nameOrEmail, password)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("VALID"))
			log.Println("SUCCESS")
	
		} else {
			log.Println("ERROR: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	
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

		ctx := r.Context()
	
		err := changeName.Service.ChangeNameFunction(ctx, currentName, newName)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("VALID"))
			log.Println("SUCCESS")
	
		} else {
			log.Println("ERROR: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

		ctx := r.Context()
	
		err := changeEmail.Service.ChangeEmailFunction(ctx, currentEmail, newEmail, confirmNewEmail, password)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("VALID"))
			log.Println("SUCCESS")
	
		} else {
			log.Println("ERROR: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

		ctx := r.Context()
	
		err, token := requestHandler.Service.RequestFunction(ctx, email)
		if err == nil {
			log.Println("SUCCESS")
			http.Redirect(w, r, "/valid?token="+token, http.StatusSeeOther)
			return
			
		} else {
			log.Println("ERROR: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}
}


func (resetPasswordHandler *ResetPasswordHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.SetFlags(log.Lshortfile)

	if r.Method == http.MethodPost {

		currentPassword := r.FormValue("currentPassword")
		newPassword := r.FormValue("newPassword")
		confirmPassword := r.FormValue("confirmNewPassword")

		ctx := r.Context()

		err := resetPasswordHandler.Service.ResetPasswordFunction(ctx, currentPassword, newPassword, confirmPassword)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("VALID"))
			log.Println("SUCCESS")

		} else {
			log.Println("ERROR: ", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

		ctx := r.Context()

		err := deleteAccountHandler.Service.DeleteAccountFunction(ctx, email, password)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("VALID"))
			log.Println("SUCCESS")

		} else {
			log.Println("ERROR")
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
	}
}


func (validToken *ValidTokenHandler) ValidTokenHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "text/html; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	log.SetFlags(log.Lshortfile)

	ctx := r.Context()

	err := validToken.Service.ValidTokenFunction(ctx)
	token := r.URL.Query().Get("token")
	if err == nil && token != "" {
		w.WriteHeader(http.StatusOK)
		tmpl.Execute(w, nil)
		log.Println("SUCCESS")

	} else {
		log.Println("ERROR: ", err)
		http.Error(w, "ERROR", http.StatusBadRequest)
		return
	}
}
