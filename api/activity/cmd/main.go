package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	_ "github.com/pd120424d/mountain-service/api/activity/cmd/docs"
	"github.com/pd120424d/mountain-service/api/activity/internal/handler"
	"github.com/pd120424d/mountain-service/api/activity/internal/middleware"
	"github.com/pd120424d/mountain-service/api/activity/internal/model"
	"github.com/pd120424d/mountain-service/api/activity/internal/publisher"
	"github.com/pd120424d/mountain-service/api/activity/internal/repositories"
	"github.com/pd120424d/mountain-service/api/activity/internal/service"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	globConf "github.com/pd120424d/mountain-service/api/shared/config"
	"github.com/pd120424d/mountain-service/api/shared/firestorex/googleadapter"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/server"

	s2semployee "github.com/pd120424d/mountain-service/api/shared/s2s/employee"
	s2surgency "github.com/pd120424d/mountain-service/api/shared/s2s/urgency"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"google.golang.org/api/option"

	// Import contracts for Swagger documentation
	_ "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	_ "github.com/pd120424d/mountain-service/api/contracts/common/v1"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// @title API Сервис за Активности
// @version 1.0
// @description Сервис за праћење активности и ревизију у систему за управљање ургентним ситуацијама

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /api/v1

// @securityDefinitions.oauth2.password OAuth2Password
// @tokenUrl /api/v1/oauth/token
// @scope.read Grants read access
// @scope.write Grants write access

// @security OAuth2Password

func main() {
	svcName := globConf.ActivityServiceName
	ctx, _ := utils.EnsureRequestID(context.Background())
	log, err := utils.NewLogger(svcName)
	if err != nil {
		panic(fmt.Sprintf("Failed to create logger: %v", err))
	}
	log.WithContext(ctx).Info("Starting Activity service")
	defer utils.TimeOperation(log, "ActivityService.main")()

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

	var client *pubsub.Client
	var err error
	credsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
	if credsPath == "" {
		credsPath = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	}
	useADC := false
	if credsPath != "" {
		if info, statErr := os.Stat(credsPath); statErr != nil || info.Size() == 0 {
			log.Warnf("Credentials path set but file missing/empty (path=%s). Falling back to ADC.", credsPath)
			useADC = true
		}
	}
	if credsPath != "" && !useADC {
		log.Infof("Initializing Pub/Sub client with credentials file: %s", credsPath)
		client, err = pubsub.NewClient(context.Background(), projectID, option.WithCredentialsFile(credsPath))
	} else {
		log.Info("Initializing Pub/Sub client using Application Default Credentials (ADC)")
		client, err = pubsub.NewClient(context.Background(), projectID)
	}
	if err != nil {
		log.Errorf("Failed to create Pub/Sub client: %v", err)
		return
	}

	repo := repositories.NewOutboxRepository(log, db)
	topic := os.Getenv("PUBSUB_TOPIC")
	if topic == "" {
		topic = "activity-events"
	}

	// Configure publisher interval and batch size via env (defaults: 10s, 100)
	intervalSec := 10
	if v := os.Getenv("OUTBOX_PUBLISH_INTERVAL_SECONDS"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			intervalSec = i
		}
	}
	batchSize := 100
	if v := os.Getenv("OUTBOX_PUBLISH_BATCH_SIZE"); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			batchSize = i
		}
	}

	pub := publisher.New(log, repo, client, publisher.Config{TopicName: topic, Interval: time.Duration(intervalSec) * time.Second, BatchSize: batchSize})
	ctx, _ := context.WithCancel(context.Background())
	pub.Start(ctx)
}

func setupRoutes(log utils.Logger, r *gin.Engine, db *gorm.DB) {
	log.Info("Setting up custom activity routes")

	// Initialize repositories and services
	activityRepo := repositories.NewActivityRepository(log, db)

	serviceAuth := auth.NewServiceAuth(auth.ServiceAuthConfig{Secret: os.Getenv("SERVICE_AUTH_SECRET"), ServiceName: "activity-service", TokenTTL: time.Hour})

	urgencyBaseURL := os.Getenv("URGENCY_SERVICE_URL")
	if urgencyBaseURL == "" {
		urgencyBaseURL = "http://urgency-service:8083"
	}
	urgencyClient := s2surgency.New(s2surgency.Config{BaseURL: urgencyBaseURL, ServiceAuth: serviceAuth, Logger: log, Timeout: 30 * time.Second})

	employeeBaseURL := os.Getenv("EMPLOYEE_SERVICE_URL")
	if employeeBaseURL == "" {
		employeeBaseURL = "http://employee-service:8082"
	}
	employeeClient := s2semployee.New(s2semployee.Config{BaseURL: employeeBaseURL, ServiceAuth: serviceAuth, Logger: log, Timeout: 30 * time.Second})

	activitySvc := service.NewActivityServiceWithDeps(log, activityRepo, urgencyClient, employeeClient)

	var flagSvc service.FeatureFlagService

	// Initialize Firestore service if env vars present
	var readModel service.FirestoreService
	projectID := os.Getenv("FIREBASE_PROJECT_ID")
	if projectID == "" {
		projectID = os.Getenv("GOOGLE_CLOUD_PROJECT")
	}
	if projectID != "" {
		credsPath := os.Getenv("FIREBASE_CREDENTIALS_PATH")
		if credsPath == "" {
			credsPath = os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
		}
		var fsClient *firestore.Client
		var err error
		useADC := false
		if credsPath != "" {
			if info, statErr := os.Stat(credsPath); statErr != nil || info.Size() == 0 {
				log.Warnf("Credentials path set but file missing/empty (path=%s). Falling back to ADC.", credsPath)
				useADC = true
			}
		}
		if credsPath != "" && !useADC {
			log.Infof("Initializing Firestore client with credentials file: %s", credsPath)
			fsClient, err = firestore.NewClient(context.Background(), projectID, option.WithCredentialsFile(credsPath))
		} else {
			log.Info("Initializing Firestore client using Application Default Credentials (ADC)")
			fsClient, err = firestore.NewClient(context.Background(), projectID)
		}
		if err != nil {
			log.Errorf("Failed to create Firestore client, continuing without read model: %v", err)
		} else {
			adapter := googleadapter.NewClientAdapter(fsClient)
			readModel = service.NewFirebaseReadService(adapter, log)
			flagSvc = service.NewFirestoreFeatureFlag(adapter, log)
			log.Info("Activity read-model (Firestore) enabled")
		}
	} else {
		log.Warn("Firestore project ID not set, fetching activities will fallback to SQL DB!")
	}

	// Fallback to in-memory flag service if Firestore not available
	if flagSvc == nil {
		flagSvc = service.NewInMemoryFeatureFlag(log, false)
	}

	// Load feature flag configuration
	defaultSource := os.Getenv("ACTIVITY_DATA_SOURCE")
	if defaultSource == "" {
		defaultSource = "firestore"
	}
	adminCanToggle := os.Getenv("ADMIN_CAN_TOGGLE_ACTIVITY_SOURCE")
	if adminCanToggle == "" {
		adminCanToggle = "true"
	}
	log.Infof("Activity data source configuration: default=%s, admin_can_toggle=%s", defaultSource, adminCanToggle)

	activityHandler := handler.NewActivityHandler(log, activitySvc, readModel, urgencyClient, defaultSource, adminCanToggle == "true")

	activityHandler.SetFeatureFlagService(flagSvc)

	// Setup JWT secret
	jwtSecret := server.SetupJWTSecret(log)
	_ = jwtSecret // JWT secret is set up but not used directly here

	redisAddr := os.Getenv(globConf.REDIS_ADDR)
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	blacklistConfig := auth.TokenBlacklistConfig{RedisAddr: redisAddr, RedisDB: 0}
	tokenBlacklist := auth.NewTokenBlacklist(blacklistConfig)
	if err := tokenBlacklist.TestConnection(); err != nil {
		log.Fatalf("Failed to connect to Redis token blacklist: %v. Redis is required for secure token invalidation.", err)
	}
	log.Info("Successfully initialized Redis token blacklist")

	authMiddleware := auth.AuthMiddleware(log, tokenBlacklist)
	adminToggleMiddleware := middleware.AdminToggleMiddleware(middleware.AdminToggleConfig{
		Logger:         log,
		AdminCanToggle: adminCanToggle == "true",
	})

	authorized := r.Group("/api/v1").Use(authMiddleware, adminToggleMiddleware)
	{
		authorized.POST("/activities", activityHandler.CreateActivity)
		authorized.GET("/activities", activityHandler.ListActivities)
		authorized.GET("/activities/counts", activityHandler.GetActivityCounts)
		authorized.GET("/activities/:id", activityHandler.GetActivity)
		authorized.DELETE("/activities/:id", activityHandler.DeleteActivity)
	}

	// Service-to-service internal routes (hidden from Swagger)
	serviceGroup := r.Group("/api/v1/service").Use(auth.NewServiceAuthMiddleware(serviceAuth))
	{
		serviceGroup.POST("/activities", activityHandler.CreateActivity)
		serviceGroup.GET("/activities", activityHandler.ListActivities)
	}

	// Admin-only routes
	admin := r.Group("/api/v1/admin").Use(auth.AdminMiddleware(log, tokenBlacklist))
	{
		admin.POST("/activities/batch", activityHandler.AddActivitiesBatch)
		admin.DELETE("/activities/reset", activityHandler.ResetAllData)
		admin.GET("/feature-flags/activity-source", activityHandler.GetActivitySourceFlag)
		admin.PUT("/feature-flags/activity-source", activityHandler.SetActivitySourceFlag)
	}
}
