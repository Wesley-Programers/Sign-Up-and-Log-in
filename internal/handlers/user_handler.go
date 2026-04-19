package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	
	"ShieldAuth-API/internal/response"
	"ShieldAuth-API/internal/middleware"
	"ShieldAuth-API/internal/service"
)


type ChangeNameHandler struct {
	Service *service.ChangeName
}
type ChangeEmailHandler struct {
	Service *service.ChangeEmail
}
type ResetPasswordHandler struct {
	Service *service.ResetPassword
}
type DeleteAccountHandler struct {
	Service *service.DeleteAccount
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


func (changeName *ChangeNameHandler) ChangeNameHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {

		var req ChangeNameRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid request", err)
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

		response.Json(w, http.StatusOK, map[string]string{"message": "success"})

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
			response.Error(w, http.StatusBadRequest, "Invalid request", err)
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

		response.Json(w, http.StatusOK, map[string]string{"message": "success"})

	} else {
		log.Println("ERROR")
		http.Error(w, "Error", http.StatusMethodNotAllowed)
		return
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
