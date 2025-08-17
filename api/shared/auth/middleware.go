package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// NewServiceAuthMiddleware creates a middleware that validates service-to-service JWT tokens
func NewServiceAuthMiddleware(serviceAuth ServiceAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := serviceAuth.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// Store service name in context for later use
		c.Set("service_name", claims.ServiceName)
		c.Next()
	}
}

// OptionalServiceAuthMiddleware creates a middleware that optionally validates service tokens
// This is useful for endpoints that can be called by both users and services
func OptionalServiceAuthMiddleware(serviceAuth ServiceAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// Try to validate as service token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) == 2 && tokenParts[0] == "Bearer" {
			token := tokenParts[1]
			if claims, err := serviceAuth.ValidateToken(token); err == nil {
				c.Set("service_name", claims.ServiceName)
				c.Set("is_service_request", true)
			}
		}

		c.Next()
	}
}
