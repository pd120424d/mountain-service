package main

import (
	"fmt"

	"github.com/pd120424d/mountain-service/api/shared/auth"
	globConf "github.com/pd120424d/mountain-service/api/shared/config"
	"github.com/pd120424d/mountain-service/api/shared/server"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	_ "github.com/pd120424d/mountain-service/api/urgency/cmd/docs"
	"github.com/pd120424d/mountain-service/api/urgency/internal"
	internalConfig "github.com/pd120424d/mountain-service/api/urgency/internal/config"
	"github.com/pd120424d/mountain-service/api/urgency/internal/model"
	"github.com/pd120424d/mountain-service/api/urgency/internal/repositories"
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

	// Setup server configuration
	serverConfig := server.ServerConfig{
		ServiceName: svcName,
		Port:        globConf.UrgencyServicePort,
		DatabaseConfig: server.GetDatabaseConfigWithDefaults(
			[]interface{}{&model.Urgency{}, &model.Notification{}},
			globConf.UrgencyDBName,
		),
		CORSConfig: server.DefaultCORSConfig(),
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
	log.Info("Setting up custom urgency routes")

	// Initialize repositories
	urgencyRepo := repositories.NewUrgencyRepository(log, db)
	notificationRepo := repositories.NewNotificationRepository(log, db)

	// Initialize service clients using defaults-aware loader (env vars override)
	serviceConfig := internalConfig.LoadServiceConfig()

	serviceClients, err := internalConfig.InitializeServiceClients(serviceConfig, log)
	if err != nil {
		log.Fatalf("Failed to initialize service clients: %v", err)
	}

	// Initialize service with all dependencies
	urgencySvc := internal.NewUrgencyService(log, urgencyRepo, notificationRepo, serviceClients.EmployeeClient)
	urgencyHandler := internal.NewUrgencyHandler(log, urgencySvc)

	// Public routes (no authentication required) - registering a new urgency
	r.POST("/api/v1/urgencies", urgencyHandler.CreateUrgency)

	// Protected routes (authentication required)
	authorized := r.Group("/api/v1").Use(auth.AuthMiddleware(log, nil))
	{
		authorized.GET("/urgencies", urgencyHandler.ListUrgencies)
		authorized.GET("/urgencies/:id", urgencyHandler.GetUrgency)
		authorized.PUT("/urgencies/:id", urgencyHandler.UpdateUrgency)
		authorized.DELETE("/urgencies/:id", urgencyHandler.DeleteUrgency)
		authorized.POST("/urgencies/:id/assign", urgencyHandler.AssignUrgency)
		authorized.DELETE("/urgencies/:id/assign", urgencyHandler.UnassignUrgency)
	}

	// Admin-only routes
	admin := r.Group("/api/v1/admin").Use(auth.AdminMiddleware(log, nil))
	{
		admin.DELETE("/urgencies/reset", urgencyHandler.ResetAllData)
	}

	// Custom swagger route (if using custom SwaggerURL)
	// r.GET("/urgency-swagger.json", func(c *gin.Context) {
	//     c.File("/docs/swagger.json")
	// })
}
