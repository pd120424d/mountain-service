package auth

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

// AuthMiddleware creates a middleware that validates user JWT tokens
func AuthMiddleware(log utils.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Error("failed to validate JWT: Missing or invalid Authorization header")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			log.Errorf("failed to validate JWT: %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Store claims in context
		ctx.Set("employeeID", claims.ID)
		ctx.Set("role", claims.Role)

		log.Info("JWT validation successful")

		ctx.Next()
	}
}

// AdminMiddleware creates a middleware that validates admin JWT tokens
func AdminMiddleware(log utils.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Error("failed to validate JWT: Missing or invalid Authorization header")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			log.Errorf("failed to validate JWT: %v", err)
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		if claims.Role != "Administrator" {
			log.Error("Access denied: Administrator role required")
			ctx.JSON(http.StatusForbidden, gin.H{"error": "Administrator access required"})
			ctx.Abort()
			return
		}

		ctx.Set("employeeID", claims.ID)
		ctx.Set("role", claims.Role)

		log.Info("Admin access granted")
		ctx.Next()
	}
}

// EmployeeData interface for basic auth middleware
type EmployeeData interface {
	GetID() uint
	GetPassword() string
	GetRole() string
}

// BasicAuthMiddleware creates a middleware that validates basic authentication
func BasicAuthMiddleware(log utils.Logger, employeeRepo interface {
	GetEmployeeByUsername(username string) (EmployeeData, error)
}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Basic ") {
			log.Error("failed to validate Basic Auth: Missing or invalid Authorization header")
			ctx.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		// Decode base64 credentials
		encodedCredentials := strings.TrimPrefix(authHeader, "Basic ")
		decodedBytes, err := base64.StdEncoding.DecodeString(encodedCredentials)
		if err != nil {
			log.Errorf("failed to decode basic auth credentials: %v", err)
			ctx.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			ctx.Abort()
			return
		}

		credentials := string(decodedBytes)
		parts := strings.SplitN(credentials, ":", 2)
		if len(parts) != 2 {
			log.Error("invalid basic auth format")
			ctx.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials format"})
			ctx.Abort()
			return
		}

		username := parts[0]
		password := parts[1]

		// Check if it's admin login
		if IsAdminLogin(username) {
			if !ValidateAdminPassword(password) {
				log.Error("Invalid admin password")
				ctx.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
				ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
				ctx.Abort()
				return
			}

			// Set admin context
			ctx.Set("employeeID", uint(0)) // Admin has ID 0
			ctx.Set("role", "Administrator")
			log.Info("Admin authenticated via Basic Auth")
			ctx.Next()
			return
		}

		// Validate regular employee
		employee, err := employeeRepo.GetEmployeeByUsername(username)
		if err != nil {
			log.Errorf("failed to retrieve employee: %v", err)
			ctx.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			ctx.Abort()
			return
		}

		if !CheckPassword(employee.GetPassword(), password) {
			log.Error("failed to verify password")
			ctx.Header("WWW-Authenticate", "Basic realm=\"Restricted\"")
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
			ctx.Abort()
			return
		}

		// Store employee info in context
		ctx.Set("employeeID", employee.GetID())
		ctx.Set("role", employee.GetRole())

		log.Info("Employee authenticated via Basic Auth")
		ctx.Next()
	}
}
