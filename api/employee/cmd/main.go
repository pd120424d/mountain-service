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

	_ "github.com/pd120424d/mountain-service/api/employee/cmd/docs"
	"github.com/pd120424d/mountain-service/api/employee/internal/handler"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	globConf "github.com/pd120424d/mountain-service/api/shared/config"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	// Import contracts for Swagger documentation
	_ "github.com/pd120424d/mountain-service/api/contracts/common/v1"
	_ "github.com/pd120424d/mountain-service/api/contracts/employee/v1"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/handlers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API Сервис за Запослене
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
	svcName := globConf.EmployeeServiceName
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
	log.Info("Starting Employee Service")

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = globConf.EmployeeDBName
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
		Models: []interface{}{&model.Employee{}, &model.Shift{}, &model.EmployeeShift{}},
	}
	db := globConf.InitDb(log, svcName, dbConfig)

	employeeRepo := repositories.NewEmployeeRepository(log, db)
	shiftsRepo := repositories.NewShiftRepository(log, db)
	employeeHandler := handler.NewEmployeeHandler(log, employeeRepo, shiftsRepo)

	r := gin.Default()

	r.Use(log.RequestLogger())

	setupRoutes(log, r, employeeHandler)

	corsHandler := setupCORS(log, r)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", globConf.EmployeeServicePort),
		Handler: corsHandler,
	}

	// Run server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	log.Infof("Starting Employee Service on port %s", globConf.EmployeeServicePort)

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

func setupCORS(log utils.Logger, r *gin.Engine) http.Handler {
	log.Info("Setting up CORS...")
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(gin.Recovery())
	r.GET("/swagger/*any", func(c *gin.Context) {
		log.Infof("Swagger request: %s %s from %s", c.Request.Method, c.Request.URL.Path, c.ClientIP())
		ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL("/swagger.json"),
		)(c)
	})
	r.GET("/swagger.json", func(c *gin.Context) {
		c.File("/docs/swagger.json")
	})

	// CORS setup
	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})

	corsOriginsEnv := os.Getenv("CORS_ALLOWED_ORIGINS") // e.g., "http://localhost:4200 for local development
	corsOrigins := strings.Split(corsOriginsEnv, ",")
	origins := handlers.AllowedOrigins(corsOrigins)

	log.Infof("Allowed CORS origins: %s", os.Getenv("CORS_ALLOWED_ORIGINS"))

	log.Info("CORS setup finished")

	// Wrap the router with CORS middleware
	return handlers.CORS(origins, headers, methods)(r)
}

func setupRoutes(log utils.Logger, r *gin.Engine, employeeHandler handler.EmployeeHandler) {
	r.POST("/api/v1/employees", employeeHandler.RegisterEmployee)
	r.POST("/api/v1/login", employeeHandler.LoginEmployee)
	r.POST("/api/v1/oauth/token", employeeHandler.OAuth2Token)
	authorized := r.Group("/api/v1").Use(auth.AuthMiddleware(log))
	{
		authorized.GET("/employees", employeeHandler.ListEmployees)
		authorized.DELETE("/employees/:id", employeeHandler.DeleteEmployee)
		authorized.POST("/employees/:id/shifts", employeeHandler.AssignShift)
		authorized.PUT("/employees/:id", employeeHandler.UpdateEmployee)
		authorized.GET("/employees/:id/shifts", employeeHandler.GetShifts)
		authorized.GET("/shifts/availability", employeeHandler.GetShiftsAvailability)
		authorized.DELETE("/employees/:id/shifts", employeeHandler.RemoveShift)

		// Service-to-service endpoints (for now using regular auth, will add service auth later)
		authorized.GET("/employees/on-call", employeeHandler.GetOnCallEmployees)
		authorized.GET("/employees/:id/active-emergencies", employeeHandler.CheckActiveEmergencies)
	}

	// Admin-only routes
	admin := r.Group("/api/v1/admin").Use(auth.AdminMiddleware(log))
	{
		admin.DELETE("/reset", employeeHandler.ResetAllData)
	}

	public := r.Group("/api/v1")
	{
		public.GET("/health", func(c *gin.Context) {
			log.Info("Health endpoint hit")
			c.JSON(200, gin.H{"message": "Service is healthy", "service": "employee"})
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
