package handlers

import (
	"log"
	"net/http"
	"time"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("PANIC: Recovered from fatal error: %v", err)
				http.Error(w, "SERVER ERROR", http.StatusInternalServerError)
			}
		}()

		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[%s] %s | Duration: %v", r.Method, r.URL.Path, time.Since(start),)
	})
}
