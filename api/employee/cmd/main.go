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

	"github.com/gorilla/handlers"

	_ "github.com/pd120424d/mountain-service/api/employee/cmd/docs"
	"github.com/pd120424d/mountain-service/api/employee/config"
	"github.com/pd120424d/mountain-service/api/employee/internal/auth"
	"github.com/pd120424d/mountain-service/api/employee/internal/handler"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API Сервис за Запослене
// @version 1.0
// @description Ово је пример API сервиса за запослене.
// @termsOfService http://example.com/terms/

// @contact.name Подршка за API
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8082
// @BasePath /api/v1
func main() {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "staging"
	}

	log, err := utils.NewLogger("employee-service")
	if err != nil {
		fmt.Errorf("failed to create logger: %v", err)
	}
	defer func(log utils.Logger) {
		err := log.Sync()
		if err != nil {
			log.Fatalf("failed to sync logger: %v", err)
		}
	}(log)

	dbUser, err := readSecret(os.Getenv("DB_USER_FILE"))
	if err != nil {
		log.Fatalf("Failed to read DB_USER: %v", err)
	}

	dbPassword, err := readSecret(os.Getenv("DB_PASSWORD_FILE"))
	if err != nil {
		log.Fatalf("Failed to read DB_PASSWORD: %v", err)
	}

	log.Infof("Connecting to database at %s:%s as user %s", dbHost, dbPort, dbUser)
	dbStringEmployee := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	// Create the employee_service database if it doesn't exist

	db := config.GetEmployeeDB(log, dbStringEmployee)

	// Auto migrate the model
	err = db.AutoMigrate(&model.Employee{})
	if err != nil {
		log.Fatalf("failed to migrate employee model: %v", err)
	}

	employeeRepo := repositories.NewEmployeeRepository(log, db)
	shiftsRepo := repositories.NewShiftRepository(log, db)
	employeeHandler := handler.NewEmployeeHandler(log, employeeRepo, shiftsRepo)

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(gin.Recovery())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("/swagger.json"),
	))
	r.GET("/swagger.json", func(c *gin.Context) {
		c.File("/docs/swagger.json")
	})

	r.POST("/api/v1/employees", employeeHandler.RegisterEmployee)
	r.POST("/api/v1/login", employeeHandler.LoginEmployee)
	authorized := r.Group("/api/v1").Use(auth.AuthMiddleware())
	{
		authorized.GET("/employees", employeeHandler.ListEmployees)
		authorized.DELETE("/employees/:id", employeeHandler.DeleteEmployee)
		authorized.POST("/empoyees/{id}/shifts", employeeHandler.AssignShift)
		authorized.GET("/employees/{id}/shifts", employeeHandler.GetShifts)
		authorized.GET("/shifts/availability", employeeHandler.GetShiftsAvailability)
	}

	// CORS setup
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	corsOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS") // e.g., "http://localhost:4200 for local development
	corsOrigins := strings.Split(corsOriginsEnv, ",")
	origins := handlers.AllowedOrigins(corsOrigins)

	log.Infof("Allowed CORS origins: %s", os.Getenv("CORS_ALLOWED_ORIGINS"))

	// Wrap the router with CORS middleware
	corsHandler := handlers.CORS(origins, headers, methods)(r)

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

	log.Infof("Starting Employee Service on port %s", config.ServerPort)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down Employee Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Employee Service forced to shutdown: %v", err)
	}

	log.Info("Employee Service exiting")
}

func readSecret(filePath string) (string, error) {
	secret, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(secret), nil
}
