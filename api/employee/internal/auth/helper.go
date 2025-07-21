package auth

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type EmployeeClaims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(employeeID uint, role string) (string, error) {
	claims := EmployeeClaims{
		ID:   employeeID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // Token expires in 24h
		},
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func IsAdminLogin(username string) bool {
	return strings.ToLower(username) == "admin"
}

func ValidateAdminPassword(password string) bool {
	adminPasswordFile := os.Getenv("ADMIN_PASSWORD_FILE")
	if adminPasswordFile == "" {
		adminPassword := os.Getenv("ADMIN_PASSWORD")
		if adminPassword == "" {
			return false
		}
		return password == adminPassword
	}

	adminPasswordBytes, err := os.ReadFile(adminPasswordFile)
	if err != nil {
		return false
	}

	adminPassword := strings.TrimSpace(string(adminPasswordBytes))
	return password == adminPassword
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
