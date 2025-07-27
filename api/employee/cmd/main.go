package main

import (
	"fmt"

	_ "github.com/pd120424d/mountain-service/api/employee/cmd/docs"
	"github.com/pd120424d/mountain-service/api/employee/internal/handler"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	globConf "github.com/pd120424d/mountain-service/api/shared/config"
	"github.com/pd120424d/mountain-service/api/shared/server"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	// Import contracts for Swagger documentation
	_ "github.com/pd120424d/mountain-service/api/contracts/common/v1"
	_ "github.com/pd120424d/mountain-service/api/contracts/employee/v1"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	// Setup server configuration
	serverConfig := server.ServerConfig{
		ServiceName: svcName,
		Port:        globConf.EmployeeServicePort,
		DatabaseConfig: server.GetDatabaseConfigWithDefaults(
			[]interface{}{&model.Employee{}, &model.Shift{}, &model.EmployeeShift{}},
			globConf.EmployeeDBName,
		),
		CORSConfig: server.CORSConfig{
			AllowedOrigins:     []string{"http://localhost:4200"},
			AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowedHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-Requested-With"},
			ExposeHeaders:      []string{"Content-Length"},
			AllowCredentials:   true,
			UseGorillaHandlers: true,
		},
		RouteConfig: server.RouteConfig{
			ServiceName: svcName,
		},
		SetupCustomRoutes: func(log utils.Logger, r *gin.Engine, db *gorm.DB) {
			setupRoutes(log, r, db)
		},
	}

	// Initialize and run server
	if err := server.InitializeServer(log, serverConfig); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
}

func setupRoutes(log utils.Logger, r *gin.Engine, db *gorm.DB) {
	log.Info("Setting up custom employee routes")

	// Initialize repositories and handler
	employeeRepo := repositories.NewEmployeeRepository(log, db)
	shiftsRepo := repositories.NewShiftRepository(log, db)
	employeeHandler := handler.NewEmployeeHandler(log, employeeRepo, shiftsRepo)

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
}
