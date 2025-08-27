package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/pubsub"
	_ "github.com/pd120424d/mountain-service/api/activity/cmd/docs"
	"github.com/pd120424d/mountain-service/api/activity/internal/handler"
	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	"github.com/pd120424d/mountain-service/api/activity/internal/publisher"
	"github.com/pd120424d/mountain-service/api/activity/internal/repositories"
	"github.com/pd120424d/mountain-service/api/activity/internal/service"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	globConf "github.com/pd120424d/mountain-service/api/shared/config"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/server"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	// Import contracts for Swagger documentation
	_ "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	_ "github.com/pd120424d/mountain-service/api/contracts/common/v1"

	"github.com/gin-gonic/gin"
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

// @host localhost:8084
// @BasePath /api/v1

// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl /api/v1/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access

// @security OAuth2Password

func main() {
	svcName := globConf.ActivityServiceName
	log, err := utils.NewLogger(svcName)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}

	// Setup server configuration
	serverConfig := server.ServerConfig{
		ServiceName: svcName,
		Port:        globConf.ActivityServicePort,
		DatabaseConfig: server.GetDatabaseConfigWithDefaults(
			[]interface{}{&model.Activity{}, &models.OutboxEvent{}},
			globConf.ActivityDBName,
		),
		CORSConfig: server.DefaultCORSConfig(),
		RouteConfig: server.RouteConfig{
			ServiceName: svcName,
		},
		SetupCustomRoutes: func(log utils.Logger, r *gin.Engine, db *gorm.DB) {
			setupRoutes(log, r, db)
			startPublisherIfConfigured(log, db)
		},
	}

	// Initialize and run server
	if err := server.InitializeServer(log, serverConfig); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}
}

func startPublisherIfConfigured(log utils.Logger, db *gorm.DB) {
	// Build Pub/Sub client if GOOGLE_APPLICATION_CREDENTIALS/FIREBASE creds are available
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if projectID == "" {
		log.Warn("Pub/Sub publisher disabled: no project ID in env")
		return
	}

	client, err := pubsub.NewClient(context.Background(), projectID)
	if err != nil {
		log.Errorf("Failed to create Pub/Sub client: %v", err)
		return
	}

	repo := repositories.NewOutboxRepository(log, db)
	topic := os.Getenv("PUBSUB_TOPIC")
	if topic == "" {
		topic = "activity-events"
	}

	pub := publisher.New(log, repo, client, publisher.Config{TopicName: topic, Interval: 10 * time.Second})
	ctx, _ := context.WithCancel(context.Background())
	pub.Start(ctx)
}

func setupRoutes(log utils.Logger, r *gin.Engine, db *gorm.DB) {
	log.Info("Setting up custom activity routes")

	// Initialize repositories and services
	activityRepo := repositories.NewActivityRepository(log, db)
	activitySvc := service.NewActivityService(log, activityRepo)
	activityHandler := handler.NewActivityHandler(log, activitySvc)

	// Setup JWT secret
	jwtSecret := server.SetupJWTSecret(log)
	_ = jwtSecret // JWT secret is set up but not used directly here

	// For activity service, we don't have user logout, so no blacklist usage required here.
	// Pass a no-op blacklist implementation.
	noopBlacklist := auth.NewTokenBlacklist(auth.TokenBlacklistConfig{RedisAddr: "localhost:6379", RedisDB: 0})
	_ = noopBlacklist

	authMiddleware := auth.AuthMiddleware(log, nil)

	authorized := r.Group("/api/v1").Use(authMiddleware)
	{
		authorized.POST("/activities", activityHandler.CreateActivity)
		authorized.GET("/activities", activityHandler.ListActivities)
		authorized.GET("/activities/stats", activityHandler.GetActivityStats)
		authorized.GET("/activities/:id", activityHandler.GetActivity)
		authorized.DELETE("/activities/:id", activityHandler.DeleteActivity)
	}

	// Admin-only routes
	admin := r.Group("/api/v1/admin").Use(auth.AdminMiddleware(log, nil))
	{
		admin.DELETE("/activities/reset", activityHandler.ResetAllData)
	}
}
