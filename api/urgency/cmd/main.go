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
	globConf "github.com/pd120424d/mountain-service/api/shared/config"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	_ "github.com/pd120424d/mountain-service/api/urgency/cmd/docs"
	"github.com/pd120424d/mountain-service/api/urgency/internal"
	internalConfig "github.com/pd120424d/mountain-service/api/urgency/internal/config"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/pd120424d/mountain-service/api/urgency/internal/repositories"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
	svcName := globConf.UrgencyServiceName
	log, err := utils.NewLogger(svcName)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	defer func(log utils.Logger) {
		err := log.Sync()
		if err != nil {
			log.Fatalf("failed to sync logger: %v", err)
		}
	}(log)
	log.Info("Starting Urgency Service")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = globConf.UrgencyDBName
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		log.Info("APP_ENV is not set, defaulting to staging")
		env = "staging"
	}

	dbConfig := globConf.DatabaseConfig{
		Host:   dbHost,
		Port:   dbPort,
		Name:   dbName,
		Models: []interface{}{&model.Urgency{}, &model.EmergencyAssignment{}, &model.Notification{}},
	}
	db := globConf.InitDb(log, svcName, dbConfig)

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
		Addr:    fmt.Sprintf(":%v", globConf.UrgencyServicePort),
		Handler: corsHandler,
	}

	// Run server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof("Starting Urgency Service on port %s", globConf.UrgencyServicePort)

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
