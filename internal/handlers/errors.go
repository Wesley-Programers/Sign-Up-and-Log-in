package handlers

import (
	"errors"
	"net/http"
	
	"ShieldAuth-API/internal/response"
	"ShieldAuth-API/internal/domain"
)

func MapServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailIsTheSame):
		response.Error(w, http.StatusConflict, "The new email is the same as the current one", err)

	case errors.Is(err, domain.ErrEmailMismatch):
		response.Error(w, http.StatusBadRequest, "The new email and its confirmation are not the same", err)

	case errors.Is(err, domain.ErrInvalidData):
		response.Error(w, http.StatusBadRequest, "Invalid input", err)

	default:
		response.Error(w, http.StatusInternalServerError, "Internal server error", err)
	}
}