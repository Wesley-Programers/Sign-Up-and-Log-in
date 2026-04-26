package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func EmailKeyFunc(r *http.Request) string {
	if r.Body == nil {
		return "invalid_request"
	}

	bodyBytes, err := io.ReadAll(io.LimitReader(r.Body, 1024))
	if err != nil {
		return "invalid_request"
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	var payload struct {
		Email string `json:"email"`
	}

	if err := json.Unmarshal(bodyBytes, &payload); err != nil || payload.Email == "" {
		return "malformed_json_attempt"
	}

	return "email:" + payload.Email
}