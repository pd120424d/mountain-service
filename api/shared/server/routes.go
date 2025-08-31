package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/swaggo/swag"
)

type RouteConfig struct {
	ServiceName string
}

func SetupHealthEndpoint(log utils.Logger, r *gin.Engine, serviceName string) {
	log.Info("Setting up health endpoint")

	public := r.Group("/api/v1")
	{
		public.GET("/health", func(c *gin.Context) {
			log.Info("Health endpoint hit")
			c.JSON(200, gin.H{
				"message": "Service is healthy",
				"service": serviceName,
			})
		})
	}
}

func SetupSwaggerEndpoints(log utils.Logger, r *gin.Engine, config RouteConfig) {
	log.Info("Setting up swagger endpoints")

	// Setup swagger.json endpoint - serve dynamically generated swagger spec with correct host
	r.GET("/swagger.json", func(c *gin.Context) {
		log.Infof("Serving swagger.json for %s", config.ServiceName)
		c.Header("Content-Type", "application/json")

		// Get the swagger spec from swag registry (this includes runtime host override)
		spec := swag.GetSwagger("swagger")
		if spec != nil {
			// ReadDoc() returns a JSON string, so we need to send it as raw data
			c.Data(200, "application/json", []byte(spec.ReadDoc()))
		} else {
			// Fallback to static file if spec not found
			c.File("/docs/swagger.json")
		}
	})
}

// SetupJWTSecret reads JWT secret from file or environment variable
func SetupJWTSecret(log utils.Logger) string {
	jwtSecret, err := ReadSecret("/run/secrets/jwt_secret")
	if err != nil {
		log.Warnf("Failed to read JWT secret from file, using environment variable: %v", err)
		jwtSecret = os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Fatal("JWT_SECRET environment variable is required")
		}
	}
	return jwtSecret
}

func SetupCommonRoutes(log utils.Logger, r *gin.Engine, serviceName string) {
	SetupHealthEndpoint(log, r, serviceName)
	SetupSwaggerEndpoints(log, r, RouteConfig{
		ServiceName: serviceName,
	})
}
