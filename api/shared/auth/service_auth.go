package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ServiceClaims struct {
	ServiceName string `json:"service"`
	jwt.RegisteredClaims
}

type ServiceAuthConfig struct {
	Secret      string
	ServiceName string
	TokenTTL    time.Duration
}

//go:generate mockgen -destination=service_auth_gomock.go -package=auth -source=service_auth.go ServiceAuth -typed
type ServiceAuth interface {
	GenerateToken() (string, error)
	ValidateToken(tokenString string) (*ServiceClaims, error)
	GetAuthHeader() (string, error)
}

type serviceAuth struct {
	config ServiceAuthConfig
}

func NewServiceAuth(config ServiceAuthConfig) ServiceAuth {
	if config.TokenTTL == 0 {
		config.TokenTTL = time.Hour // Default to 1 hour
	}
	return &serviceAuth{config: config}
}

func (sa *serviceAuth) GenerateToken() (string, error) {
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
func (sa *serviceAuth) ValidateToken(tokenString string) (*ServiceClaims, error) {
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

func (sa *serviceAuth) GetAuthHeader() (string, error) {
	token, err := sa.GenerateToken()
	if err != nil {
		return "", err
	}
	return "Bearer " + token, nil
}

type EmployeeClaims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(employeeID uint, role string) (string, error) {
	now := time.Now()
	claims := EmployeeClaims{
		ID:   employeeID,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(), // Unique token ID for blacklisting
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)), // Token expires in 24h
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ValidateJWT validates a user JWT token and checks blacklist
func ValidateJWT(tokenString string, blacklist TokenBlacklist) (*EmployeeClaims, error) {
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

	// Check blacklist if provided
	if blacklist != nil && claims.RegisteredClaims.ID != "" {
		ctx := context.Background()
		isBlacklisted, err := blacklist.IsTokenBlacklisted(ctx, claims.RegisteredClaims.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to check token blacklist: %w", err)
		}
		if isBlacklisted {
			return nil, errors.New("token has been revoked")
		}
	}

	return claims, nil
}

func GenerateAdminJWT() (string, error) {
	now := time.Now()
	claims := EmployeeClaims{
		ID:   0,
		Role: "Administrator",
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        uuid.New().String(), // Unique token ID for blacklisting
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24)),
			NotBefore: jwt.NewNumericDate(now),
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
