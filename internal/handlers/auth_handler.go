package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"text/template"
	"time"

	"ShieldAuth-API/internal/domain"
	"ShieldAuth-API/internal/response"
	"ShieldAuth-API/internal/security"
	"ShieldAuth-API/internal/service"
	"ShieldAuth-API/internal/ui"
)


type RegisterHandler struct {
	Service *service.Register
}
type LoginHandler struct {
	Service *service.VerifyLogin
	Limiter *security.RedisLimiter
}
type RequestHandler struct {
	Service *service.Request
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
func NewRequestHandler(service *service.Request) *RequestHandler {
	return &RequestHandler{
		Service: service,
	}
}
func NewValidTokenHandler(service *service.ValidToken) *ValidTokenHandler {
	return &ValidTokenHandler{
		Service: service,
	}
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
				response.Error(w, http.StatusConflict, "Email already exists", err)

			case errors.Is(err, domain.ErrInvalidData):
				response.Error(w, http.StatusBadRequest, "Invalid input", err)

			default:
				log.Printf("Unexpected error: %v", err)
				response.Error(w, http.StatusInternalServerError, "Internal server error", err)
			}
			return
		}
		response.Json(w, http.StatusOK, map[string]string{"message": "success"})
	
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
			response.Error(w, http.StatusBadRequest, "Invalid credentials", err)
			return
		}

		key := "login-attempt:" + req.NameOrEmail
		allowed, err := login.Limiter.CheckLimit(r.Context(), key, 5, 10*time.Minute)
		if err != nil {
			response.Error(w, http.StatusInternalServerError, "Internal server error", err)
			return
		}

		if !allowed {
			response.Error(w, http.StatusTooManyRequests, "Too many requests, try again later", err)
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
			response.Error(w, http.StatusInternalServerError, "ERROR TRYING TO CREATE A TOKEN", err)
			return
		}

		log.Println("SUCCESS")
		response.Json(w, http.StatusOK, map[string]string{"token": tokenJwtString})

	} else {
		log.Println("ERROR")
		http.Error(w, "ERROR", http.StatusMethodNotAllowed)
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

