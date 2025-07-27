package server

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouteConfig struct {
	ServiceName    string
	SwaggerURL     string // Custom swagger URL (optional)
	SwaggerDocPath string // Path to swagger.json file (default: "/docs/swagger.json")
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

	if config.SwaggerDocPath == "" {
		config.SwaggerDocPath = "/docs/swagger.json"
	}

	// Setup swagger UI endpoint
	if config.SwaggerURL != "" {
		// Custom swagger URL (like urgency service)
		r.GET("/swagger/*any", func(c *gin.Context) {
			log.Infof("Swagger request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
			ginSwagger.WrapHandler(swaggerFiles.Handler,
				ginSwagger.URL(config.SwaggerURL),
			)(c)
		})
	} else {
		// Standard swagger setup
		r.GET("/swagger/*any", func(c *gin.Context) {
			log.Infof("Swagger request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
			ginSwagger.WrapHandler(swaggerFiles.Handler,
				ginSwagger.URL("/swagger.json"),
			)(c)
		})
	}

	// Setup swagger.json endpoint
	r.GET("/swagger.json", func(c *gin.Context) {
		c.File(config.SwaggerDocPath)
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
