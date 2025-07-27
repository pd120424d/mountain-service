package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type CORSConfig struct {
	AllowedOrigins     []string
	AllowedMethods     []string
	AllowedHeaders     []string
	ExposeHeaders      []string
	AllowCredentials   bool
	UseGorillaHandlers bool
}

func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins:     []string{"http://localhost:3000", "http://localhost:4200", "http://localhost:8080", "http://localhost:8081", "http://localhost:8082", "http://localhost:8083", "http://localhost:8084"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization", "X-Requested-With"},
		ExposeHeaders:      []string{"Content-Length"},
		AllowCredentials:   true,
		UseGorillaHandlers: false,
	}
}

func SetupCORS(log utils.Logger, r *gin.Engine, config CORSConfig) http.Handler {
	log.Info("Setting up CORS...")

	if len(config.AllowedOrigins) == 0 {
		corsOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
		if corsOriginsEnv != "" {
			config.AllowedOrigins = strings.Split(corsOriginsEnv, ",")
		} else {
			config.AllowedOrigins = DefaultCORSConfig().AllowedOrigins
		}
	}

	if len(config.AllowedMethods) == 0 {
		config.AllowedMethods = DefaultCORSConfig().AllowedMethods
	}
	if len(config.AllowedHeaders) == 0 {
		config.AllowedHeaders = DefaultCORSConfig().AllowedHeaders
	}
	if len(config.ExposeHeaders) == 0 {
		config.ExposeHeaders = DefaultCORSConfig().ExposeHeaders
	}

	log.Infof("Allowed CORS origins: %v", config.AllowedOrigins)

	if config.UseGorillaHandlers {
		r.Use(cors.New(cors.Config{
			AllowOrigins:     config.AllowedOrigins,
			AllowMethods:     config.AllowedMethods,
			AllowHeaders:     config.AllowedHeaders,
			ExposeHeaders:    config.ExposeHeaders,
			AllowCredentials: config.AllowCredentials,
		}))

		r.Use(gin.Recovery())

		headers := handlers.AllowedHeaders(config.AllowedHeaders)
		methods := handlers.AllowedMethods(config.AllowedMethods)
		origins := handlers.AllowedOrigins(config.AllowedOrigins)

		log.Info("CORS setup finished")
		return handlers.CORS(origins, headers, methods)(r)
	} else {
		corsConfig := cors.Config{
			AllowOrigins:     config.AllowedOrigins,
			AllowMethods:     config.AllowedMethods,
			AllowHeaders:     config.AllowedHeaders,
			ExposeHeaders:    config.ExposeHeaders,
			AllowCredentials: config.AllowCredentials,
		}

		r.Use(cors.New(corsConfig))
		log.Info("CORS setup finished")
		return r
	}
}
