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

	"github.com/gin-contrib/cors"
	"github.com/gorilla/handlers"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	_ "github.com/pd120424d/mountain-service/api/urgency/cmd/docs"
	"github.com/pd120424d/mountain-service/api/urgency/config"
	"github.com/pd120424d/mountain-service/api/urgency/internal"
	internalConfig "github.com/pd120424d/mountain-service/api/urgency/internal/config"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/pd120424d/mountain-service/api/urgency/internal/repositories"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	// Import contracts for Swagger documentation
	_ "github.com/pd120424d/mountain-service/api/contracts/common/v1"
	_ "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"

	"github.com/gin-gonic/gin"
)

// @title API Сервис за Ургентне ситуације
// @version 1.0

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl /api/v1/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access

// @security OAuth2Password

// @host
// @BasePath /api/v1
func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "staging"
	}

	log, err := utils.NewLogger("urgency-service")
	if err != nil {
		panic("failed to create logger:" + err.Error())
	}
	defer func(log utils.Logger) {
		err := log.Sync()
		if err != nil {
			log.Fatalf("failed to sync logger: %v", err)
		}
	}(log)

	db := initDb(log, dbHost, dbPort, dbName)

	// Initialize repositories
	urgencyRepo := repositories.NewUrgencyRepository(log, db)
	assignmentRepo := repositories.NewAssignmentRepository(log, db)
	notificationRepo := repositories.NewNotificationRepository(log, db)

	// Initialize service clients
	serviceConfig := internalConfig.ServiceConfig{
		EmployeeServiceURL: os.Getenv("EMPLOYEE_SERVICE_URL"),
		ServiceAuthSecret:  os.Getenv("SERVICE_AUTH_SECRET"),
		ServiceName:        "urgency-service",
	}

	serviceClients, err := internalConfig.InitializeServiceClients(serviceConfig, log)
	if err != nil {
		log.Fatalf("Failed to initialize service clients: %v", err)
	}

	// Initialize service with all dependencies
	urgencySvc := internal.NewUrgencyService(log, urgencyRepo, assignmentRepo, notificationRepo, serviceClients.EmployeeClient)
	urgencyHandler := internal.NewUrgencyHandler(log, urgencySvc)

	r := gin.Default()

	r.Use(log.RequestLogger())

	setupRoutes(log, r, urgencyHandler)

	corsHandler := setupCORS(log, r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", config.ServerPort),
		Handler: corsHandler,
	}

	// Run server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof("Starting Urgency Service on port %s", config.ServerPort)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down Urgency Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Urgency Service forced to shutdown: %v", err)
	}

	log.Info("Urgency Service exiting")
}

func initDb(log utils.Logger, dbHost, dbPort, dbName string) *gorm.DB {
	log.Info("Setting up database...")

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		userFile := os.Getenv("URGENCY_DB_USER_FILE")
		if userFile != "" {
			var err error
			dbUser, err = readSecret(userFile)
			if err != nil {
				log.Fatalf("Failed to read URGENCY_DB_USER from file %s: %v", userFile, err)
			}
		} else {
			log.Fatalf("Neither DB_USER environment variable nor URGENCY_DB_USER_FILE is set")
		}
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		passwordFile := os.Getenv("URGENCY_DB_PASSWORD_FILE")
		if passwordFile != "" {
			var err error
			dbPassword, err = readSecret(passwordFile)
			if err != nil {
				log.Fatalf("Failed to read URGENCY_DB_PASSWORD from file %s: %v", passwordFile, err)
			}
		} else {
			log.Fatalf("Neither DB_PASSWORD environment variable nor URGENCY_DB_PASSWORD_FILE is set")
		}
	}

	log.Infof("Connecting to database at %s:%s as user %s", dbHost, dbPort, dbUser)
	dbStringUrgency := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Create the urgency_service database if it doesn't exist
	db := config.GetUrgencyDB(log, dbStringUrgency)

	// Auto migrate the models
	err := db.AutoMigrate(&model.Urgency{}, &model.EmergencyAssignment{}, &model.Notification{})
	if err != nil {
		log.Fatalf("failed to migrate urgency models: %v", err)
	}
	log.Info("Successfully migrated urgency models")

	log.Info("Database setup finished")
	return db
}

func setupCORS(log utils.Logger, r *gin.Engine) http.Handler {
	log.Info("Setting up CORS...")

	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	r.Use(cors.New(corsConfig))

	r.Use(gin.Recovery())
	r.GET("/swagger/*any", func(c *gin.Context) {
		log.Infof("Swagger request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL("/urgency-swagger.json"),
		)(c)
	})
	r.GET("/swagger.json", func(c *gin.Context) {
		c.File("/docs/swagger.json")
	})

	// CORS setup
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	corsOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	corsOrigins := strings.Split(corsOriginsEnv, ",")
	origins := handlers.AllowedOrigins(corsOrigins)

	log.Infof("Allowed CORS origins: %s", os.Getenv("CORS_ALLOWED_ORIGINS"))

	log.Info("CORS setup finished")

	return handlers.CORS(origins, headers, methods)(r)
}

func setupRoutes(log utils.Logger, r *gin.Engine, urgencyHandler internal.UrgencyHandler) {
	authorized := r.Group("/api/v1").Use(auth.AuthMiddleware(log))
	{
		authorized.POST("/urgencies", urgencyHandler.CreateUrgency)
		authorized.GET("/urgencies", urgencyHandler.ListUrgencies)
		authorized.GET("/urgencies/:id", urgencyHandler.GetUrgency)
		authorized.PUT("/urgencies/:id", urgencyHandler.UpdateUrgency)
		authorized.DELETE("/urgencies/:id", urgencyHandler.DeleteUrgency)
	}

	// Admin-only routes
	admin := r.Group("/api/v1/admin").Use(auth.AdminMiddleware(log))
	{
		admin.DELETE("/urgencies/reset", urgencyHandler.ResetAllData)
	}

	public := r.Group("/api/v1")
	{
		public.GET("/health", func(c *gin.Context) {
			log.Info("Health endpoint hit")
			c.JSON(200, gin.H{"message": "Service is healthy", "service": "urgency"})
		})
	}
}

func readSecret(filePath string) (string, error) {
	secret, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(secret)), nil
}
