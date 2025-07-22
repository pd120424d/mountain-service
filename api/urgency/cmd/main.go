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
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/config"
	"github.com/pd120424d/mountain-service/api/urgency/internal/auth"
	"github.com/pd120424d/mountain-service/api/urgency/internal/handler"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/pd120424d/mountain-service/api/urgency/internal/repositories"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

// @title API Сервис за Хитности
// @version 1.0

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

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

	urgencyRepo := repositories.NewUrgencyRepository(log, db)
	urgencyHandler := handler.NewUrgencyHandler(log, urgencyRepo)

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
	dbUser, err := readSecret(os.Getenv("URGENCY_DB_USER_FILE"))
	if err != nil {
		log.Fatalf("Failed to read URGENCY_DB_USER: %v", err)
	}

	dbPassword, err := readSecret(os.Getenv("URGENCY_DB_PASSWORD_FILE"))
	if err != nil {
		log.Fatalf("Failed to read URGENCY_DB_PASSWORD: %v", err)
	}

	log.Infof("Connecting to database at %s:%s as user %s", dbHost, dbPort, dbUser)
	dbStringUrgency := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Create the urgency_service database if it doesn't exist
	db := config.GetUrgencyDB(log, dbStringUrgency)

	// Auto migrate the model
	err = db.AutoMigrate(&model.Urgency{})
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

	// Add Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"*"}),
	)(r)

	log.Info("CORS setup finished")
	return corsHandler
}

func setupRoutes(log utils.Logger, r *gin.Engine, urgencyHandler handler.UrgencyHandler) {
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

	// Health check
	r.GET("/ping", func(c *gin.Context) {
		log.Info("Ping route hit")
		c.JSON(200, gin.H{"message": "pong"})
	})
}

func readSecret(filePath string) (string, error) {
	secret, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(secret)), nil
}
