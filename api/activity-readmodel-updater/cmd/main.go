package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/option"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/repositories"
	"github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/service"
	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/models"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type Config struct {
	DatabaseURL string

	FirebaseProjectID       string
	FirebaseCredentialsPath string

	PubSubTopic        string
	PubSubSubscription string

	OutboxPollIntervalSeconds int

	HealthPort int

	LogLevel string

	Version string
	GitSHA  string
}

func loadConfig() *Config {
	return &Config{
		DatabaseURL:               getEnvOrDefault("DATABASE_URL", "postgres://user:password@localhost:5432/activities?sslmode=disable"),
		FirebaseProjectID:         getEnvOrDefault("FIREBASE_PROJECT_ID", "your-project-id"),
		FirebaseCredentialsPath:   getEnvOrDefault("FIREBASE_CREDENTIALS_PATH", ""),
		PubSubTopic:               getEnvOrDefault("PUBSUB_TOPIC", "activity-events"),
		PubSubSubscription:        getEnvOrDefault("PUBSUB_SUBSCRIPTION", "activity-events-sub"),
		OutboxPollIntervalSeconds: getEnvAsIntOrDefault("OUTBOX_POLL_INTERVAL_SECONDS", 10),
		HealthPort:                getEnvAsIntOrDefault("HEALTH_PORT", 8090),
		LogLevel:                  getEnvOrDefault("LOG_LEVEL", "info"),
		Version:                   getEnvOrDefault("VERSION", "dev"),
		GitSHA:                    getEnvOrDefault("GIT_SHA", "unknown"),
	}
}

func main() {
	// Load configuration
	cfg := loadConfig()

	// Initialize logger
	logger, err := utils.NewLogger("activity-readmodel-updater")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	logger.Infof("Starting Activity Read Model Updater version=%s git_sha=%s", cfg.Version, cfg.GitSHA)

	// Initialize database connection
	db, err := initDatabase(cfg.DatabaseURL, logger)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Firebase Firestore
	firestoreClient, err := initFirestore(cfg.FirebaseCredentialsPath, cfg.FirebaseProjectID)
	if err != nil {
		logger.Fatalf("Failed to initialize Firestore: %v", err)
	}
	defer firestoreClient.Close()

	// Initialize Pub/Sub
	pubsubClient, err := initPubSub(cfg.FirebaseCredentialsPath, cfg.FirebaseProjectID)
	if err != nil {
		logger.Fatalf("Failed to initialize Pub/Sub: %v", err)
	}
	defer pubsubClient.Close()

	// Initialize services
	firebaseService := service.NewFirebaseService(firestoreClient, logger)
	outboxRepo := repositories.NewOutboxRepository(logger, db)

	// Start health check server
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"healthy","service":"activity-readmodel-updater"}`))
		})

		healthAddr := fmt.Sprintf(":%d", cfg.HealthPort)
		logger.Infof("Starting health check server on %s", healthAddr)
		if err := http.ListenAndServe(healthAddr, nil); err != nil {
			logger.Errorf("Health check server error: %v", err)
		}
	}()

	// Start outbox publisher (polls outbox table and publishes to Pub/Sub)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		ticker := time.NewTicker(time.Duration(cfg.OutboxPollIntervalSeconds) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if err := processOutboxEvents(ctx, outboxRepo, pubsubClient, cfg.PubSubTopic, logger); err != nil {
					logger.Errorf("Failed to process outbox events: %v", err)
				}
			}
		}
	}()

	// Start Pub/Sub subscriber (receives events and updates Firestore)
	go func() {
		subscription := pubsubClient.Subscription(cfg.PubSubSubscription)
		logger.Infof("Starting Pub/Sub subscriber: subscription=%s", cfg.PubSubSubscription)

		err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			if err := handleActivityEvent(ctx, msg, firebaseService, logger); err != nil {
				logger.Errorf("Failed to handle activity event: error=%v, message_id=%s", err, msg.ID)
				msg.Nack()
			} else {
				msg.Ack()
			}
		})

		if err != nil {
			logger.Errorf("Pub/Sub subscriber error: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	logger.Infof("Activity Read Model Updater started successfully")
	<-sigChan

	logger.Infof("Shutting down Activity Read Model Updater...")
	cancel()

	// Give some time for graceful shutdown
	time.Sleep(5 * time.Second)
	logger.Infof("Activity Read Model Updater stopped")
}

func initDatabase(databaseURL string, logger utils.Logger) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Infof("Database connection established")
	return db, nil
}

func initFirestore(credentialsPath, projectID string) (*firestore.Client, error) {
	ctx := context.Background()

	var client *firestore.Client
	var err error

	if credentialsPath != "" {
		client, err = firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	} else {
		client, err = firestore.NewClient(ctx, projectID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Firestore client: %w", err)
	}

	return client, nil
}

func initPubSub(credentialsPath, projectID string) (*pubsub.Client, error) {
	ctx := context.Background()

	var client *pubsub.Client
	var err error

	if credentialsPath != "" {
		client, err = pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	} else {
		client, err = pubsub.NewClient(ctx, projectID)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Pub/Sub client: %w", err)
	}

	return client, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntOrDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// processOutboxEvents processes unpublished events from the outbox and publishes them to Pub/Sub
func processOutboxEvents(ctx context.Context, outboxRepo repositories.OutboxRepository, pubsubClient *pubsub.Client, topicName string, logger utils.Logger) error {
	logger.Infof("Processing outbox events")

	// Get unpublished events
	events, err := outboxRepo.GetUnpublishedEvents(100) // Process up to 100 events at a time
	if err != nil {
		logger.Errorf("Failed to get unpublished events: %v", err)
		return fmt.Errorf("failed to get unpublished events: %w", err)
	}

	if len(events) == 0 {
		logger.Infof("No unpublished events found")
		return nil
	}

	logger.Infof("Found %d unpublished events to process", len(events))

	// Get the topic
	topic := pubsubClient.Topic(topicName)
	defer topic.Stop()

	// Process each event
	successCount := 0
	for _, event := range events {
		if err := publishEvent(ctx, topic, event, logger); err != nil {
			logger.Errorf("Failed to publish event: event_id=%d, error=%v", event.ID, err)
			continue
		}

		if err := outboxRepo.MarkAsPublished(event.ID); err != nil {
			logger.Errorf("Failed to mark event as published: event_id=%d, error=%v", event.ID, err)
			continue
		}

		successCount++
	}

	logger.Infof("Successfully processed %d out of %d outbox events", successCount, len(events))
	return nil
}

func publishEvent(ctx context.Context, topic *pubsub.Topic, event *models.OutboxEvent, logger utils.Logger) error {
	logger.Infof("Publishing event to Pub/Sub: event_id=%d, type=%s", event.ID, event.EventType)

	messageBytes, err := json.Marshal(event)
	if err != nil {
		logger.Errorf("Failed to marshal event: event_id=%d, error=%v", event.ID, err)
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	pubsubMessage := &pubsub.Message{
		Data: messageBytes,
		Attributes: map[string]string{
			"event_type":   event.EventType,
			"aggregate_id": event.AggregateID,
		},
	}

	result := topic.Publish(ctx, pubsubMessage)
	messageID, err := result.Get(ctx)
	if err != nil {
		logger.Errorf("Failed to publish message to Pub/Sub: event_id=%d, error=%v", event.ID, err)
		return fmt.Errorf("failed to publish message to Pub/Sub: %w", err)
	}

	logger.Infof("Event published to Pub/Sub successfully: event_id=%d, message_id=%s", event.ID, messageID)
	return nil
}

// handleActivityEvent processes an activity event message from Pub/Sub
func handleActivityEvent(ctx context.Context, msg *pubsub.Message, firebaseService service.FirebaseService, logger utils.Logger) error {
	logger.Infof("Received activity event: message_id=%s, publish_time=%v", msg.ID, msg.PublishTime)

	var event activityV1.OutboxEvent
	if err := json.Unmarshal(msg.Data, &event); err != nil {
		logger.Errorf("Failed to unmarshal event message: error=%v, message_id=%s", err, msg.ID)
		return fmt.Errorf("failed to unmarshal event message: %w", err)
	}

	logger.Infof("Processing activity event: event_type=%s, aggregate_id=%s", event.EventType, event.AggregateID)

	// Extract the ActivityEvent from the outbox event
	activityEvent, err := event.GetEventData()
	if err != nil {
		logger.Errorf("Failed to extract activity event data: error=%v", err)
		return fmt.Errorf("failed to extract activity event data: %w", err)
	}

	if err := firebaseService.SyncActivity(ctx, *activityEvent); err != nil {
		logger.Errorf("Failed to sync activity to Firebase: activity_id=%d, error=%v", activityEvent.ActivityID, err)
		return fmt.Errorf("failed to sync activity to Firebase: %w", err)
	}

	logger.Infof("Activity synced to Firebase successfully: activity_id=%d", activityEvent.ActivityID)
	return nil
}
