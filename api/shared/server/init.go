package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"gorm.io/gorm"
)

type ServerConfig struct {
	ServiceName string
	Port        string

	DatabaseConfig DatabaseConfig

	CORSConfig CORSConfig

	RouteConfig RouteConfig

	SetupCustomRoutes func(log utils.Logger, r *gin.Engine, db *gorm.DB)
}

// InitializeServer initializes a complete server with database, CORS, routes, and graceful shutdown
func InitializeServer(log utils.Logger, config ServerConfig) error {
	log.Infof("Initializing server: %s", config.ServiceName)

	db := InitDb(log, config.ServiceName, config.DatabaseConfig)

	r := gin.Default()
	// Ensure every request carries a request ID and echo it back in the response
	r.Use(utils.RequestIDMiddleware())
	r.Use(log.RequestLogger())

	// Apply fresh-read window from client (if any)
	r.Use(utils.FreshReadWindowMiddleware())

	SetupHealthEndpoint(log, r, config.ServiceName)
	SetupSwaggerEndpoints(log, r, config.RouteConfig)

	if config.SetupCustomRoutes != nil {
		config.SetupCustomRoutes(log, r, db)
	}

	corsHandler := SetupCORS(log, r, config.CORSConfig)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.Port),
		Handler: corsHandler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof("%s started on port %s", config.ServiceName, config.Port)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infof("Shutting down %s...", config.ServiceName)

	// Give outstanding requests a deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("%s forced to shutdown: %v", config.ServiceName, err)
	}

	log.Infof("%s exited", config.ServiceName)
	return nil
}

func GetDatabaseConfig(models []interface{}) DatabaseConfig {
	return DatabaseConfig{
		Host:   os.Getenv("DB_HOST"),
		Port:   os.Getenv("DB_PORT"),
		Name:   os.Getenv("DB_NAME"),
		Models: models,
	}
}

func GetDatabaseConfigWithDefaults(models []interface{}, defaultDBName string) DatabaseConfig {
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = defaultDBName
	}

	return DatabaseConfig{
		Host:   os.Getenv("DB_HOST"),
		Port:   os.Getenv("DB_PORT"),
		Name:   dbName,
		Models: models,
	}
}
