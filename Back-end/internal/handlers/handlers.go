package handlers

import (
	"encoding/json"
	"errors"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"index/Back-end/internal/domain"
	"index/Back-end/internal/middleware"
	"index/Back-end/internal/security"
	"index/Back-end/internal/service"
	"index/Back-end/internal/ui"
	"index/Back-end/internal/web"
)

type RegisterHandler struct {
	Service *service.Register
}

type LoginHandler struct {
	Service *service.VerifyLogin
	Limiter *security.RedisLimiter
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


type ChangeNameRequest struct {
	CurrentName string `json:"currentName" validate:"required,name"`
	NewName string `json:"newName" validate:"required,name"`
	ConfirmNewName string `json:"confirmNewName" validate:"eqfield=NewName"`
}

type ChangeEmailRequest struct {
	CurrentEmail string `json:"currentEmail" validate:"required,email"`
	NewEmail string `json:"newEmail" validate:"required,email"`
	ConfirmNewEmail string `json:"confirmNewEmail" validate:"eqfield=NewEmail"`
	Password string `json:"password" validate:"required"`
}

type LoginRequest struct {
	NameOrEmail string `json:"nameOrEmail" validate:"required"`
	Password string `json:"password" validate:"required"`
}

var tmpl = template.Must(template.ParseFS(ui.Files, "templates/reset.html"))

func (handler *RegisterHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {

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
			case errors.Is(err, domain.ErrEmailAlreadyExists):
				web.Error(w, http.StatusConflict, "Email already exists", err)

			case errors.Is(err, domain.ErrInvalidData):
				web.Error(w, http.StatusBadRequest, "Invalid input", err)

			default:
				log.Printf("Unexpected error: %v", err)
				web.Error(w, http.StatusInternalServerError, "Internal server error", err)
			}
			return
		}
		web.Json(w, http.StatusOK, map[string]string{"message": "success"})
	
	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
	}
}


func (login *LoginHandler) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {

		var req LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			web.Error(w, http.StatusBadRequest, "Invalid credentials", err)
			return
		}

		key := "login-attempt:" + req.NameOrEmail
		allowed, err := login.Limiter.CheckLimit(r.Context(), key, 5, 10*time.Minute)
		if err != nil {
			web.Error(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}

		if !allowed {
			web.Error(w, http.StatusTooManyRequests, "Too many requests, try again later", err)
			return
		}

		input := service.LoginData{
			Name: req.NameOrEmail,
			Email: req.NameOrEmail,
			Password: req.Password,
		}
	
		err, id := login.Service.VerifyLoginFunction(r.Context(), input)
		if err != nil {
			MapServiceError(w, err)
			return
		}

		login.Limiter.ResetLimit(r.Context(), "login-attempt:"+req.NameOrEmail)

		tokenJwtString, err := service.TokenJWT(id)
		if err != nil {
			web.Error(w, http.StatusInternalServerError, "ERROR TRYING TO CREATE A TOKEN", err)
			return
		}

		log.Println("SUCCESS")
		web.Json(w, http.StatusOK, map[string]string{"token": tokenJwtString})

	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
		return
	}
}


func (changeName *ChangeNameHandler) ChangeNameHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {

		var req ChangeNameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			web.Error(w, http.StatusBadRequest, "Invalid request", err)
			return
		}

		idString, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		idContext, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		input := service.ChangeNameData{
			ID: idContext,
			CurrentName: req.CurrentName,
			NewName: req.NewName,
			ConfirmNewName: req.ConfirmNewName,
		}

		if err := changeName.Service.ChangeNameFunction(r.Context(), input); err != nil {
			MapServiceError(w, err)
			return
		}

		web.Json(w, http.StatusOK, map[string]string{"message": "success"})

	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
		return
	}
}


func (changeEmail *ChangeEmailHandler) ChangeEmailHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {

		var req ChangeEmailRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			web.Error(w, http.StatusBadRequest, "Invalid request", err)
			return
		}
	
		idString, ok := r.Context().Value(middleware.UserIDKey).(string)
		if !ok {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		idContext, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		input := service.ChangeEmailData{
			ID: idContext,
			CurrentEmail: req.CurrentEmail,
			NewEmail: req.NewEmail,
			ConfirmNewEmail: req.ConfirmNewEmail,
			Password: req.Password,
		}
	
		if err := changeEmail.Service.ChangeEmailFunctionTest(r.Context(), input); err != nil {
			MapServiceError(w, err)
			return
		}

		web.Json(w, http.StatusOK, map[string]string{"message": "success"})

	} else {
		log.Println("ERROR")
		http.Error(w, "Error", http.StatusMethodNotAllowed)
		return
	}
}


func (requestHandler *RequestHandler) RequestHandler(w http.ResponseWriter, r *http.Request) {

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
