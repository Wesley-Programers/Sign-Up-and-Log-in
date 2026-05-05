package handlers

import (
	"encoding/json"
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
	service service.Service
}


func NewRegisterHanlder(service *service.Register) *RegisterHandler {
	return &RegisterHandler{
		Service: service,
	}
}
func NewLoginHandler(service *service.VerifyLogin, limiter *security.RedisLimiter) *LoginHandler {
	return &LoginHandler{
		Service: service,
		Limiter: limiter,
	}
}
func NewRequestHandler(s service.Service) *RequestHandler {
	return &RequestHandler{
		service: s,
	}
}


type RegisterRequest struct {
	Name string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required"`
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

		var req RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.Error(w, http.StatusBadRequest, "Invalid credentials", err)
			return
		}

		input := service.RegisterData{
			Name: req.Name,
			Email: req.Email,
			Password: req.Password,
		}
	
		err := handler.Service.RegisterFunction(r.Context(), input)
		if err != nil {
			MapServiceError(w, err)
			return
		}
		response.Json(w, http.StatusCreated, map[string]string{"message": "success"})
	
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


func (h *RequestHandler) RequestReset(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request", err)
		return
	}

	_, err := h.service.RequestReset(r.Context(), req.Email)
	if err != nil {
		MapServiceError(w, err)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "if email exists, reset link was sent",
	})
}


func (h *RequestHandler) ValidToken(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		MapServiceError(w, domain.ErrInvalidToken)
		return
	}

	_, err := h.service.ValidToken(r.Context(), token)
	if err != nil {
		MapServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "valid",
	})
}