package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	events "github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/event"

	"cloud.google.com/go/firestore"
	"cloud.google.com/go/pubsub"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"

	"github.com/pd120424d/mountain-service/api/activity-readmodel-updater/internal/service"
	"github.com/pd120424d/mountain-service/api/shared/firestorex/googleadapter"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type Config struct {
	DatabaseURL string

	FirebaseProjectID       string
	FirebaseCredentialsPath string

	PubSubTopic        string
	PubSubSubscription string

	OutboxPollIntervalSeconds int

	// Pub/Sub subscriber parallelism
	SubscriberNumGoroutines          int
	SubscriberMaxOutstandingMessages int
	SubscriberMaxOutstandingBytes    int

	// Internal sharded dispatcher parallelism (per-activity ordering)
	ShardWorkers int
	ShardQueue   int

	HealthPort int

	LogLevel string

	Version string
	GitSHA  string
}

func loadConfig() *Config {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = buildDatabaseURLFromEnv()
	}

	// Prefer explicit FIREBASE_CREDENTIALS_PATH; fallback to GOOGLE_APPLICATION_CREDENTIALS
	credPath := getEnvOrDefault("FIREBASE_CREDENTIALS_PATH", "")
	if credPath == "" {
		credPath = getEnvOrDefault("GOOGLE_APPLICATION_CREDENTIALS", "")
	}
	return &Config{
		DatabaseURL:                      dbURL,
		FirebaseProjectID:                getEnvOrDefault("FIREBASE_PROJECT_ID", "your-project-id"),
		FirebaseCredentialsPath:          credPath,
		PubSubTopic:                      getEnvOrDefault("PUBSUB_TOPIC", "activity-events"),
		PubSubSubscription:               getEnvOrDefault("PUBSUB_SUBSCRIPTION", "activity-events-sub"),
		OutboxPollIntervalSeconds:        getEnvAsIntOrDefault("OUTBOX_POLL_INTERVAL_SECONDS", 10),
		SubscriberNumGoroutines:          getEnvAsIntOrDefault("SUBSCRIBER_NUM_GOROUTINES", 8),
		SubscriberMaxOutstandingMessages: getEnvAsIntOrDefault("SUBSCRIBER_MAX_OUTSTANDING_MESSAGES", 1000),
		SubscriberMaxOutstandingBytes:    getEnvAsIntOrDefault("SUBSCRIBER_MAX_OUTSTANDING_BYTES", 100*1024*1024),
		ShardWorkers:                     getEnvAsIntOrDefault("SHARD_WORKERS", 16),
		ShardQueue:                       getEnvAsIntOrDefault("SHARD_QUEUE", 1024),
		HealthPort:                       getEnvAsIntOrDefault("HEALTH_PORT", 8090),
		LogLevel:                         getEnvOrDefault("LOG_LEVEL", "info"),
		Version:                          getEnvOrDefault("VERSION", "dev"),
		GitSHA:                           getEnvOrDefault("GIT_SHA", "unknown"),
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
	credsSrc := "ADC"
	if cfg.FirebaseCredentialsPath != "" {
		credsSrc = fmt.Sprintf("file:%s", cfg.FirebaseCredentialsPath)
	}
	logger.Infof("Starting Activity Read Model Updater version=%s git_sha=%s project_id=%s topic=%s subscription=%s creds=%s", cfg.Version, cfg.GitSHA, cfg.FirebaseProjectID, cfg.PubSubTopic, cfg.PubSubSubscription, credsSrc)

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
	fsAdapter := googleadapter.NewClientAdapter(firestoreClient)
	firebaseService := service.NewFirebaseService(fsAdapter, logger)

	// Start health check server
	go func() {
		http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()
			// Perform a lightweight Firestore operation
			_, err := firestoreClient.Collections(ctx).Next()
			if err != nil && err != iterator.Done {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(fmt.Sprintf(`{"status":"unhealthy","service":"activity-readmodel-updater","firestore_error":"%v"}`, err)))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"healthy","service":"activity-readmodel-updater"}`))
		})

		healthAddr := fmt.Sprintf(":%d", cfg.HealthPort)
		logger.Infof("Starting health check server on %s", healthAddr)
		if err := http.ListenAndServe(healthAddr, nil); err != nil {
			logger.Errorf("Health check server error: %v", err)
		}
	}()

	// Create a cancellable context for subscriber
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start Pub/Sub subscriber
	go func() {
		subscription := pubsubClient.Subscription(cfg.PubSubSubscription)
		subscription.ReceiveSettings.NumGoroutines = cfg.SubscriberNumGoroutines
		subscription.ReceiveSettings.MaxOutstandingMessages = cfg.SubscriberMaxOutstandingMessages
		subscription.ReceiveSettings.MaxOutstandingBytes = cfg.SubscriberMaxOutstandingBytes
		logger.Infof("Starting Pub/Sub subscriber: subscription=%s", cfg.PubSubSubscription)

		// Sharded dispatcher ensures per-activity ordering, parallel across activities
		dispatcher := events.NewShardedDispatcher(firebaseService, logger, cfg.ShardWorkers, cfg.ShardQueue)
		err := subscription.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			ctx, reqID := utils.EnsureRequestID(ctx)
			reqLog := logger.WithContext(ctx)
			reqLog.Infof("Handling activity event: message_id=%s", msg.ID)

			if err := dispatcher.Process(ctx, msg); err != nil {
				reqLog.Errorf("Failed to handle activity event: error=%v, message_id=%s, request_id=%s", err, msg.ID, reqID)
				msg.Nack()
			} else {
				reqLog.Infof("Successfully handled activity event: message_id=%s", msg.ID)
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

func initFirestore(credentialsPath, projectID string) (*firestore.Client, error) {
	ctx := context.Background()

	var client *firestore.Client
	var err error

	useADC := false
	if credentialsPath != "" {
		if info, statErr := os.Stat(credentialsPath); statErr != nil || info.Size() == 0 {
			fmt.Printf("[WARN] FIREBASE_CREDENTIALS_PATH set but file missing/empty (path=%s). Falling back to ADC.\n", credentialsPath)
			useADC = true
		}
	}
	if credentialsPath != "" && !useADC {
		fmt.Printf("[INFO] Initializing Firestore with credentials file: %s\n", credentialsPath)
		client, err = firestore.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	} else {
		fmt.Println("[INFO] Initializing Firestore using Application Default Credentials (ADC)")
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

	useADC := false
	if credentialsPath != "" {
		if info, statErr := os.Stat(credentialsPath); statErr != nil || info.Size() == 0 {
			fmt.Printf("[WARN] FIREBASE_CREDENTIALS_PATH set but file missing/empty (path=%s). Falling back to ADC.\n", credentialsPath)
			useADC = true
		}
	}
	if credentialsPath != "" && !useADC {
		fmt.Printf("[INFO] Initializing Pub/Sub with credentials file: %s\n", credentialsPath)
		client, err = pubsub.NewClient(ctx, projectID, option.WithCredentialsFile(credentialsPath))
	} else {
		fmt.Println("[INFO] Initializing Pub/Sub using Application Default Credentials (ADC)")
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

// buildDatabaseURLFromEnv constructs a Postgres DSN from DB_* environment variables.
// Example: postgres://user:pass@127.0.0.1:5432/activity_service_db?sslmode=disable
func buildDatabaseURLFromEnv() string {
	user := getEnvOrDefault("DB_USER", "user")
	password := getEnvOrDefault("DB_PASSWORD", "password")
	host := getEnvOrDefault("DB_HOST", "127.0.0.1")
	port := getEnvOrDefault("DB_PORT", "5432")
	name := getEnvOrDefault("DB_NAME", "activities")
	sslmode := getEnvOrDefault("DB_SSLMODE", "disable")
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", user, password, host, port, name, sslmode)
}
