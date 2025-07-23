package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
