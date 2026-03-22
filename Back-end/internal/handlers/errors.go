package handlers

import (
	"errors"
	"net/http"
	"index/Back-end/internal/domain"
	"index/Back-end/internal/web"
)

func MapServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, domain.ErrEmailIsTheSame):
		web.Error(w, http.StatusConflict, "The new email is the same as the current one", err)

	case errors.Is(err, domain.ErrEmailMismatch):
		web.Error(w, http.StatusBadRequest, "The new email and its confirmation are not the same", err)

	case errors.Is(err, domain.ErrInvalidData):
		web.Error(w, http.StatusBadRequest, "Invalid input", err)

	default:
		web.Error(w, http.StatusInternalServerError, "Internal server error", err)
	}
}