package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthClaims struct {
	ID   uint   `json:"id"`
	Role string `json:"role"`
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			ctx.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		// Store claims in context
		ctx.Set("employeeID", claims.ID)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
