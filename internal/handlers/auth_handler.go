package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
type ValidTokenHandler struct {
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
func NewValidTokenHandler(s service.Service) *ValidTokenHandler {
	return &ValidTokenHandler{
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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

}


func (login *LoginHandler) HandlerLogin(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

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

	tokenJwtString, err := security.TokenJWT(id)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "error trying to create a token", err)
		return
	}

	response.Json(w, http.StatusOK, map[string]string{"token": tokenJwtString})
	
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

	token, err := h.service.RequestReset(r.Context(), req.Email)
	if err != nil {
		MapServiceError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"redirect": "http://127.0.0.1:8000/valid?token=" + url.QueryEscape(token),
	})
}


func (h *ValidTokenHandler) ValidToken(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	fmt.Println("token: ", token)
	if token == "" {
		MapServiceError(w, domain.ErrInvalidToken)
		return
	}

	_, err := h.service.ValidToken(r.Context(), token)
	if err != nil {
		MapServiceError(w, err)
		return
	}

	data := struct {
		Token string
	}{
		Token: token,
	}

	tmpl.Execute(w, data)
}