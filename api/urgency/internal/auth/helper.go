package auth

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type EmployeeClaims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func ValidateJWT(tokenString string) (*EmployeeClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &EmployeeClaims{}, func(token *jwt.Token) (any, error) {
		jwtSecret := []byte(os.Getenv("JWT_SECRET"))
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*EmployeeClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func GenerateAdminJWT() (string, error) {
	claims := EmployeeClaims{
		ID:   0,
		Role: "Administrator",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}
