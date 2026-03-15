package handlers

import (
	"encoding/json"
	"net/http"
	"log"
)

func Json(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func Error(w http.ResponseWriter, status int, message string, err error) {
	if err != nil {
		log.Printf("[ERROR] Status: %d | Message: %s | Internal Error: %v", status, message, err)
		Json(w, status, map[string]string{"error": message})
	}
}