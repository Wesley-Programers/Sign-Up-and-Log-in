package handlers

import (
	"errors"
	"net/http"
	
	"ShieldAuth-API/internal/response"
	"ShieldAuth-API/internal/domain"
)

func MapServiceError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, domain.ErrInvalidData), 
         errors.Is(err, domain.ErrInvalidEmailFormat):
        response.Error(w, http.StatusBadRequest, "Invalid input data", err)

    case errors.Is(err, domain.ErrEmailsDoNotMismatch):
        response.Error(w, http.StatusBadRequest, "The new email and its confirmation do not match", err)

    case errors.Is(err, domain.ErrEmailMismatch):
        response.Error(w, http.StatusBadRequest, "The provided current email does not match our records", err)

    case errors.Is(err, domain.ErrInvalidCredentials), 
         errors.Is(err, domain.ErrInvalidPassword):
        response.Error(w, http.StatusUnauthorized, "Authentication failed", err)

    case errors.Is(err, domain.ErrUserNotFound):
        response.Error(w, http.StatusNotFound, "User not found", err)

    case errors.Is(err, domain.ErrEmailAlreadyExists):
        response.Error(w, http.StatusConflict, "This email is already in use", err)

    case errors.Is(err, domain.ErrNameIsTheSame):
        response.Error(w, http.StatusConflict, "The new name is the same as the current one", err)

    case errors.Is(err, domain.ErrEmailIsTheSame):
        response.Error(w, http.StatusConflict, "The new email is the same as the current one", err)

    case errors.Is(err, domain.ErrRateLimitExceeded):
        response.Error(w, http.StatusTooManyRequests, "Too many requests, slow down", err)

    default:
        response.Error(w, http.StatusInternalServerError, "An internal server error occurred", err)
    }
}