package security

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)


func TokenJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	tokenJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secretKey := os.Getenv("JWT_KEY")
	return tokenJwt.SignedString([]byte(secretKey))
}