package middleware

import (
	"net/http"
	"context"
	"strings"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string
const UserIDKey ContextKey = "userID"

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
	
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing: %v", t.Header["alg"])
				}
				return []byte(secretKey), nil
			})
	
			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
	
			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				http.Error(w, "Invalid claims", http.StatusUnauthorized)
				return
			}
	
			sub, ok := claims["sub"]
			if !ok || sub == "" {
				http.Error(w, "Invalid subject", http.StatusUnauthorized)
				return
			}

			userID := fmt.Sprintf("%v", sub)
	
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}