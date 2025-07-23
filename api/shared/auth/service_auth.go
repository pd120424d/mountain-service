package auth

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// ServiceClaims represents the claims for service-to-service JWT tokens
type ServiceClaims struct {
	ServiceName string `json:"service"`
	jwt.RegisteredClaims
}

// ServiceAuthConfig holds the configuration for service authentication
type ServiceAuthConfig struct {
	Secret      string
	ServiceName string
	TokenTTL    time.Duration
}

// ServiceAuth handles JWT token generation and validation for service-to-service communication
type ServiceAuth struct {
	config ServiceAuthConfig
}

// NewServiceAuth creates a new ServiceAuth instance
func NewServiceAuth(config ServiceAuthConfig) *ServiceAuth {
	if config.TokenTTL == 0 {
		config.TokenTTL = time.Hour // Default to 1 hour
	}
	return &ServiceAuth{config: config}
}

// GenerateToken creates a new JWT token for service-to-service communication
func (sa *ServiceAuth) GenerateToken() (string, error) {
	now := time.Now()
	claims := ServiceClaims{
		ServiceName: sa.config.ServiceName,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(sa.config.TokenTTL)),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    sa.config.ServiceName,
			Subject:   "service-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sa.config.Secret))
}

// ValidateToken validates a JWT token and returns the service name if valid
func (sa *ServiceAuth) ValidateToken(tokenString string) (*ServiceClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ServiceClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(sa.config.Secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*ServiceClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// GetAuthHeader returns the Authorization header value for HTTP requests
func (sa *ServiceAuth) GetAuthHeader() (string, error) {
	token, err := sa.GenerateToken()
	if err != nil {
		return "", err
	}
	return "Bearer " + token, nil
}

// EmployeeClaims represents the claims for user JWT tokens
type EmployeeClaims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateJWT creates a new JWT token for user authentication
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

// ValidateJWT validates a user JWT token and returns the claims
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

// GenerateAdminJWT creates a JWT token for admin access
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

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPassword verifies a password against its hash
func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// IsAdminLogin checks if the username is for admin login
func IsAdminLogin(username string) bool {
	return strings.ToLower(username) == "admin"
}

// ValidateAdminPassword validates admin password from environment or file
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
