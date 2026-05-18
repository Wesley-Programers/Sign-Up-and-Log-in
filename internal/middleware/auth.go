package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string
const Key ContextKey = "userID"

type Claims struct {
	UserID int `json:"sub"`
	jwt.RegisteredClaims
}
type AuthContext struct {
	UserID int
	TokenHash string
}

func AuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
	
			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			if tokenString == "" || tokenString == "null" {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}

				return []byte(secretKey), nil
			},)
	
			if err != nil || !token.Valid {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			secondClaims := token.Claims.(jwt.MapClaims)
			jti, ok := secondClaims["jti"].(string)
			if !ok {
				http.Error(w, "invalid token jti", http.StatusUnauthorized)
				return
			}

			sub, ok := secondClaims["sub"].(float64)
			if !ok {
				http.Error(w, "invalid token sub", http.StatusUnauthorized)
				return
			}

			userID := int(sub)

			auth := AuthContext{
				UserID: userID,
				TokenHash: jti,
			}
	
			ctx := context.WithValue(r.Context(), Key, auth)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}