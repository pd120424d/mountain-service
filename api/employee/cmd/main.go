package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "api/employee/cmd/docs"
	"api/employee/config"
	"api/employee/internal/handler"
	"api/employee/internal/model"
	"api/employee/internal/repositories"
	"api/shared/utils"

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

	log := utils.NewLogger()
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
	shiftsRepo := repositories.NewShiftRepository(db)
	employeeHandler := handler.NewEmployeeHandler(log, employeeRepo, shiftsRepo)

	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		api.POST("/employees", employeeHandler.RegisterEmployee)
		api.GET("/employees", employeeHandler.ListEmployees)
		api.DELETE("/employees/:id", employeeHandler.DeleteEmployee)
		api.POST("/empoyees/{id}/shifts", employeeHandler.AssignShift)
		api.GET("/employees/{id}/shifts", employeeHandler.GetShifts)
		api.GET("shifts/availability", employeeHandler.GetShiftsAvailability)
	}

	// CORS setup
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:4200"})

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
