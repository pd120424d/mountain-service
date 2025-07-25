package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	_ "github.com/pd120424d/mountain-service/api/activity/cmd/docs"
	"github.com/pd120424d/mountain-service/api/activity/config"
	"github.com/pd120424d/mountain-service/api/activity/internal"
	"github.com/pd120424d/mountain-service/api/activity/internal/handler"
	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	"github.com/pd120424d/mountain-service/api/activity/internal/repositories"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	// Import contracts for Swagger documentation
	_ "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	_ "github.com/pd120424d/mountain-service/api/contracts/common/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

// @title Activity Service API
// @version 1.0
// @description Activity tracking and audit service for the Mountain Emergency Management System
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host
// @BasePath /api/v1

// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl /api/v1/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access

// @security OAuth2Password

func main() {
	log, err := utils.NewLogger("activity-service")
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	log.Info("Starting Activity Service")

	// Read environment variables
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "5432"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "mountain_service"
	}

	db := initDb(log, dbHost, dbPort, dbName)

	// Initialize repositories
	activityRepo := repositories.NewActivityRepository(log, db)

	// Initialize service
	activitySvc := internal.NewActivityService(log, activityRepo)
	activityHandler := handler.NewActivityHandler(log, activitySvc)

	r := gin.Default()

	r.Use(log.RequestLogger())

	setupRoutes(log, r, activityHandler)

	corsHandler := setupCORS(log, r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.ServerPort),
		Handler: corsHandler,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof("Activity Service started on port %v", config.ServerPort)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down Activity Service...")

	// Give outstanding requests a 30 seconds deadline for completion
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Activity Service forced to shutdown: %v", err)
	}

	log.Info("Activity Service exited")
}

func initDb(log utils.Logger, dbHost, dbPort, dbName string) *gorm.DB {
	log.Info("Setting up database...")

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		userFile := os.Getenv("ACTIVITY_DB_USER_FILE")
		if userFile != "" {
			var err error
			dbUser, err = readSecret(userFile)
			if err != nil {
				log.Fatalf("Failed to read ACTIVITY_DB_USER from file %s: %v", userFile, err)
			}
		} else {
			log.Fatalf("Neither DB_USER environment variable nor ACTIVITY_DB_USER_FILE is set")
		}
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		passwordFile := os.Getenv("ACTIVITY_DB_PASSWORD_FILE")
		if passwordFile != "" {
			var err error
			dbPassword, err = readSecret(passwordFile)
			if err != nil {
				log.Fatalf("Failed to read ACTIVITY_DB_PASSWORD from file %s: %v", passwordFile, err)
			}
		} else {
			log.Fatalf("Neither DB_PASSWORD environment variable nor ACTIVITY_DB_PASSWORD_FILE is set")
		}
	}

	log.Infof("Connecting to database at %s:%s as user %s", dbHost, dbPort, dbUser)
	dbStringActivity := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Create the activity_service database if it doesn't exist
	db := config.GetActivityDB(log, dbStringActivity)

	// Auto migrate the model
	err := db.AutoMigrate(&model.Activity{})
	if err != nil {
		log.Fatalf("failed to migrate activity models: %v", err)
	}
	log.Info("Successfully migrated activity models")

	log.Info("Database setup finished")
	return db
}

func setupCORS(log utils.Logger, r *gin.Engine) http.Handler {
	log.Info("Setting up CORS")

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:3000", "http://localhost:8080", "http://localhost:8081", "http://localhost:8082", "http://localhost:8083"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))
	return r
}

func setupRoutes(log utils.Logger, r *gin.Engine, activityHandler handler.ActivityHandler) {
	log.Info("Setting up routes")

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/swagger.json", func(c *gin.Context) {
		c.File("/docs/swagger.json")
	})

	jwtSecret, err := readSecret("/run/secrets/jwt_secret")
	if err != nil {
		log.Warnf("Failed to read JWT secret from file, using environment variable: %v", err)
		jwtSecret = os.Getenv("JWT_SECRET")
		if jwtSecret == "" {
			log.Fatal("JWT_SECRET environment variable is required")
		}
	}

	authMiddleware := auth.AuthMiddleware(log)

	public := r.Group("/api/v1")
	{
		public.GET("/health", func(c *gin.Context) {
			log.Info("Health endpoint hit")
			c.JSON(200, gin.H{"message": "Service is healthy", "service": "activity"})
		})
	}

	authorized := r.Group("/api/v1").Use(authMiddleware)
	{
		authorized.POST("/activities", activityHandler.CreateActivity)
		authorized.GET("/activities", activityHandler.ListActivities)
		authorized.GET("/activities/stats", activityHandler.GetActivityStats)
		authorized.GET("/activities/:id", activityHandler.GetActivity)
		authorized.DELETE("/activities/:id", activityHandler.DeleteActivity)
	}

	// Admin-only routes
	admin := r.Group("/api/v1/admin").Use(auth.AdminMiddleware(log))
	{
		admin.DELETE("/activities/reset", activityHandler.ResetAllData)
	}
}

func readSecret(filePath string) (string, error) {
	secret, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(secret)), nil
}
