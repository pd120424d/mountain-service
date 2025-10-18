package main

import (
	"context"
	"fmt"
	"os"
	"time"

	_ "github.com/pd120424d/mountain-service/api/employee/cmd/docs"
	"github.com/pd120424d/mountain-service/api/employee/internal/handler"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	"github.com/pd120424d/mountain-service/api/employee/internal/service"

	"github.com/pd120424d/mountain-service/api/shared/auth"
	globConf "github.com/pd120424d/mountain-service/api/shared/config"
	"github.com/pd120424d/mountain-service/api/shared/server"
	"github.com/pd120424d/mountain-service/api/shared/storage"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	// Import contracts for Swagger documentation
	_ "github.com/pd120424d/mountain-service/api/contracts/common/v1"
	_ "github.com/pd120424d/mountain-service/api/contracts/employee/v1"

	"github.com/gin-gonic/gin"
	"github.com/spf13/afero"
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

// @BasePath /api/v1

func main() {
	svcName := globConf.EmployeeServiceName
	log, err := utils.NewLogger(svcName)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}

	ctx, _ := utils.EnsureRequestID(context.Background())
	log.WithContext(ctx).Info("Starting Employee service")
	defer utils.TimeOperation(log, "EmployeeService.main")()
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
	log.Info("Setting up custom employee routes")

	// Dual-DB approach: for shifts repo, optionally use read replica if DB_READ_HOST/PORT are set
	readHost := os.Getenv("DB_READ_HOST")
	readPort := os.Getenv("DB_READ_PORT")
	var shiftsRepo repositories.ShiftRepository
	if readHost != "" && readPort != "" {
		log.Infof("DB write endpoint: %s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
		log.Infof("DB read endpoint: %s:%s", readHost, readPort)
		readCfg := server.DatabaseConfig{Host: readHost, Port: readPort, Name: os.Getenv(globConf.EmployeeDBName)}
		readDB := server.InitDb(log, globConf.EmployeeServiceName+"-read", readCfg)
		shiftsRepo = repositories.NewShiftRepositoryRW(log, db, readDB)
		log.Info("Shift repository initialized in RW mode (reads->replica, writes->primary)")
	} else {
		log.Infof("DB write endpoint: %s:%s", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"))
		log.Info("DB read endpoint: disabled (using primary)")
		shiftsRepo = repositories.NewShiftRepository(log, db)
		log.Info("Shift repository initialized in primary only mode (reads->primary, writes->primary)")
	}

	// Initialize repositories
	employeeRepo := repositories.NewEmployeeRepository(log, db)

	// Initialize Redis token blacklist
	redisAddr := os.Getenv(globConf.REDIS_ADDR)
	if redisAddr == "" {
		redisAddr = "redis:6379" // Default to local docker redis service
	}
	blacklistConfig := auth.TokenBlacklistConfig{RedisAddr: redisAddr, RedisDB: 0}
	tokenBlacklist := auth.NewTokenBlacklist(blacklistConfig)
	if err := tokenBlacklist.TestConnection(); err != nil {
		log.Fatalf("Failed to connect to Redis token blacklist: %v. Redis is required for secure token invalidation.", err)
	}
	log.Info("Successfully initialized Redis token blacklist")

	// Initialize services
	employeeService := service.NewEmployeeService(log, employeeRepo, tokenBlacklist)
	shiftService := service.NewShiftService(log, employeeRepo, shiftsRepo)

	// Initialize Azure Blob Storage service
	containerName := os.Getenv("AZURE_STORAGE_CONTAINER_NAME")
	if containerName == "" {
		containerName = "employee-profiles"
	}
	blobConfig := storage.AzureBlobConfig{
		AccountName:   os.Getenv("AZURE_STORAGE_ACCOUNT_NAME"),
		AccountKey:    os.Getenv("AZURE_STORAGE_ACCOUNT_KEY"),
		ContainerName: containerName,
	}
	log.Infof("Azure Storage Config - Account Name: %s, Container: %s, Key Present: %t",
		blobConfig.AccountName,
		blobConfig.ContainerName,
		blobConfig.AccountKey != "")

	var blobService storage.AzureBlobService
	if blobConfig.AccountName != "" && blobConfig.AccountKey != "" {
		clientWrapper, err := storage.NewAzureBlobClientWrapper(log, blobConfig)
		if err != nil {
			log.Warnf("Failed to create Azure Blob client wrapper: %v. File upload will be disabled.", err)
		} else {
			blobService, err = storage.NewAzureBlobService(log, clientWrapper)
			if err != nil {
				log.Warnf("Failed to initialize Azure Blob Storage: %v. File upload will be disabled.", err)
				blobService = nil
			} else {
				log.Info("Azure Blob Storage initialized successfully")
			}
		}
	} else {
		log.Warn("Azure Storage credentials not provided. File upload will be disabled.")
	}

	// Initialize file handler
	var fileHandler handler.FileHandler
	if blobService != nil {
		fileHandler = handler.NewFileHandler(log, blobService, employeeService)
	}

	// Initialize handler with services
	employeeHandler := handler.NewEmployeeHandler(log, afero.NewOsFs(), employeeService, shiftService)

	r.POST("/api/v1/employees", employeeHandler.RegisterEmployee)
	r.POST("/api/v1/login", employeeHandler.LoginEmployee)
	r.POST("/api/v1/oauth/token", employeeHandler.OAuth2Token)
	authorized := r.Group("/api/v1").Use(auth.AuthMiddleware(log, tokenBlacklist))
	{
		authorized.POST("/logout", employeeHandler.LogoutEmployee)
		authorized.GET("/employees", employeeHandler.ListEmployees)
		authorized.GET("/employees/:id", employeeHandler.GetEmployee)
		authorized.DELETE("/employees/:id", employeeHandler.DeleteEmployee)
		authorized.POST("/employees/:id/shifts", employeeHandler.AssignShift)
		authorized.PUT("/employees/:id", employeeHandler.UpdateEmployee)
		authorized.GET("/employees/:id/shifts", employeeHandler.GetShifts)
		authorized.GET("/employees/:id/shift-warnings", employeeHandler.GetShiftWarnings)
		authorized.GET("/shifts/availability", employeeHandler.GetShiftsAvailability)
		authorized.DELETE("/employees/:id/shifts", employeeHandler.RemoveShift)

		// Service-to-service endpoints with service authentication
		serviceAuthSecret := os.Getenv("SERVICE_AUTH_SECRET")
		if serviceAuthSecret == "" {
			log.Warn("SERVICE_AUTH_SECRET not set, service-to-service authentication may not work properly")
		}
		serviceAuth := auth.NewServiceAuth(auth.ServiceAuthConfig{
			Secret:      serviceAuthSecret,
			ServiceName: "employee-service",
			TokenTTL:    time.Hour,
		})
		serviceAuthMiddleware := auth.NewServiceAuthMiddleware(serviceAuth)

		serviceRoutes := r.Group("/api/v1").Use(serviceAuthMiddleware)
		{
			// Service-to-service: expose minimal read endpoints needed by other services
			serviceRoutes.GET("/service/employees/:id", employeeHandler.GetEmployee)
			serviceRoutes.GET("/employees/on-call", employeeHandler.GetOnCallEmployees)
			serviceRoutes.GET("/employees/:id/active-emergencies", employeeHandler.CheckActiveEmergencies)
		}

		// File upload endpoints
		if fileHandler != nil {
			authorized.POST("/employees/:id/profile-picture", fileHandler.UploadProfilePicture)
			authorized.DELETE("/employees/:id/profile-picture", fileHandler.DeleteProfilePicture)
			authorized.GET("/files/profile-picture/info", fileHandler.GetProfilePictureInfo)
		}

		// Error catalog endpoint, no need for authorization
		r.GET("/api/v1/errors/catalog", employeeHandler.GetErrorCatalog)
	}

	// Admin-only routes
	admin := r.Group("/api/v1/admin").Use(auth.AdminMiddleware(log, tokenBlacklist))
	{
		admin.DELETE("/reset", employeeHandler.ResetAllData)
		admin.GET("/shifts/availability", employeeHandler.GetAdminShiftsAvailability)
		admin.GET("/employees/:id/shift-warnings", employeeHandler.GetShiftWarnings)
		// Admin K8s ops
		admin.POST("/k8s/restart", employeeHandler.RestartDeployment)
	}
}
