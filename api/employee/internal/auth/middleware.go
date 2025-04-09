package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type AuthClaims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
}

func AuthMiddleware(log utils.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			log.Error("failed to validate JWT:Missing or invalid Authorization header")
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
