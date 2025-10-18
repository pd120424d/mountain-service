package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type AdminToggleConfig struct {
	Logger         utils.Logger
	AdminCanToggle bool
}

func AdminToggleMiddleware(config AdminToggleConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestedSource := c.GetHeader("X-Activity-Source")

		if requestedSource != "" {
			if !config.AdminCanToggle {
				config.Logger.Warn("Activity source toggle attempted but feature is disabled")
				c.JSON(http.StatusForbidden, gin.H{"error": "Activity source toggle is disabled"})
				c.Abort()
				return
			}

			role, exists := c.Get("role")
			if !exists || role != "admin" {
				config.Logger.Warnf("Non-admin user attempted to toggle activity source: role=%v", role)
				c.JSON(http.StatusForbidden, gin.H{"error": "Only admins can toggle activity source"})
				c.Abort()
				return
			}

			if requestedSource != "postgres" && requestedSource != "firestore" {
				config.Logger.Warnf("Invalid activity source requested: %s", requestedSource)
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity source. Must be 'postgres' or 'firestore'"})
				c.Abort()
				return
			}

			c.Set("activity_source_override", requestedSource)
			config.Logger.Infof("Admin toggled activity source to: %s", requestedSource)
		}

		c.Next()
	}
}

