package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"

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
	cfg := loadConfig()

	log, err := utils.NewLogger("activity-readmodel-updater")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	mainCtx, _ := utils.EnsureRequestID(context.Background())
	log = log.WithName("main").WithContext(mainCtx)
	defer utils.TimeOperation(log, "ActivityReadModelUpdaterService.main")()

	credsSrc := "ADC"
	if cfg.FirebaseCredentialsPath != "" {
		credsSrc = fmt.Sprintf("file:%s", cfg.FirebaseCredentialsPath)
	}
	log.Infof("Starting Activity Read Model Updater version=%s git_sha=%s project_id=%s topic=%s subscription=%s creds=%s", cfg.Version, cfg.GitSHA, cfg.FirebaseProjectID, cfg.PubSubTopic, cfg.PubSubSubscription, credsSrc)

	firestoreClient, err := initFirestore(mainCtx, cfg.FirebaseCredentialsPath, cfg.FirebaseProjectID)
	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	defer firestoreClient.Close()

	pubsubClient, err := initPubSub(mainCtx, cfg.FirebaseCredentialsPath, cfg.FirebaseProjectID)
	if err != nil {
		log.Fatalf("Failed to initialize Pub/Sub: %v", err)
	}
	defer pubsubClient.Close()

	fsAdapter := googleadapter.NewClientAdapter(firestoreClient)
	firebaseService := service.NewFirebaseService(fsAdapter, log)

	// Start health/ready endpoints server
	go func() {
		mux := http.NewServeMux()

		// Liveness: process is up
		mux.HandleFunc("/live", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"alive","service":"activity-readmodel-updater"}`))
		})

		// Readiness: strict dependency checks (longer timeout)
		mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()
			// Firestore connectivity
			if _, err := firestoreClient.Collections(ctx).Next(); err != nil && err != iterator.Done {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(fmt.Sprintf(`{"status":"unready","service":"activity-readmodel-updater","firestore_error":"%v"}`, err)))
				return
			}
			// Pub/Sub topic/subscription existence
			topic := pubsubClient.Topic(cfg.PubSubTopic)
			sub := pubsubClient.Subscription(cfg.PubSubSubscription)
			tExists, tErr := topic.Exists(ctx)
			sExists, sErr := sub.Exists(ctx)
			if tErr != nil || sErr != nil || !tExists || !sExists {
				w.WriteHeader(http.StatusServiceUnavailable)
				w.Write([]byte(fmt.Sprintf(`{"status":"unready","service":"activity-readmodel-updater","pubsub_topic_exists":%t,"pubsub_subscription_exists":%t,"topic_error":"%v","subscription_error":"%v"}`, tExists, sExists, tErr, sErr)))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status":"ready","service":"activity-readmodel-updater"}`))
		})

		// Backward-compatible health - lightweight summary
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
			defer cancel()
			fsOK := true
			if _, err := firestoreClient.Collections(ctx).Next(); err != nil && err != iterator.Done {
				fsOK = false
			}
			tExists, sExists := false, false
			topic := pubsubClient.Topic(cfg.PubSubTopic)
			sub := pubsubClient.Subscription(cfg.PubSubSubscription)
			if te, se := func() (error, error) {
				te := error(nil)
				se := error(nil)
				tExists, te = topic.Exists(ctx)
				sExists, se = sub.Exists(ctx)
				return te, se
			}(); te != nil || se != nil {
				log.Warnf("Failed to check Pub/Sub health: topic_exists=%t, subscription_exists=%t, topic_error=%v, subscription_error=%v", tExists, sExists, te, se)
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{"status":"ok","service":"activity-readmodel-updater","firestore_ok":%t,"pubsub_topic_exists":%t,"pubsub_subscription_exists":%t}`, fsOK, tExists, sExists)))
		})

		healthAddr := fmt.Sprintf(":%d", cfg.HealthPort)
		log.Infof("Starting health check server on %s", healthAddr)
		if err := http.ListenAndServe(healthAddr, mux); err != nil {
			log.Errorf("Health check server error: %v", err)
		}
	}()

	// Create a cancellable context for subscriber
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	ctxWithCancel, _ = utils.EnsureRequestID(ctxWithCancel)
	pubSubLog := log.WithContext(ctxWithCancel)
	defer cancel()

	// Start Pub/Sub subscriber
	go func() {
		subscription := pubsubClient.Subscription(cfg.PubSubSubscription)
		subscription.ReceiveSettings.NumGoroutines = cfg.SubscriberNumGoroutines
		subscription.ReceiveSettings.MaxOutstandingMessages = cfg.SubscriberMaxOutstandingMessages
		subscription.ReceiveSettings.MaxOutstandingBytes = cfg.SubscriberMaxOutstandingBytes

		pubSubLog.Infof("Starting Pub/Sub subscriber with config: subscription=%s num_goroutines=%d max_outstanding_messages=%d max_outstanding_bytes=%d", cfg.PubSubSubscription, cfg.SubscriberNumGoroutines, cfg.SubscriberMaxOutstandingMessages, cfg.SubscriberMaxOutstandingBytes)

		// Sharded dispatcher ensures per-activity ordering, parallel across activities
		dispatcher := events.NewShardedDispatcher(firebaseService, pubSubLog, cfg.ShardWorkers, cfg.ShardQueue)

		backoff := time.Second
		attempt := 0
		for {
			if ctxWithCancel.Err() != nil {
				pubSubLog.Infof("Subscriber context canceled; exiting receive loop")
				return
			}
			attempt++
			pubSubLog.Infof("Starting Pub/Sub subscriber receive attempt=%d subscription=%s", attempt, cfg.PubSubSubscription)
			err := subscription.Receive(ctxWithCancel, func(ctx context.Context, msg *pubsub.Message) {
				defer func() {
					if r := recover(); r != nil {
						pubSubLog.Errorf("Panic in subscriber handler: %v\nstack=%s\nmessage_id=%s", r, string(debug.Stack()), msg.ID)
						msg.Nack()
					}
				}()
				ctx, reqID := utils.EnsureRequestID(ctx)
				reqLog := pubSubLog.WithContext(ctx)
				attempt := 0
				if msg.DeliveryAttempt != nil {
					attempt = *msg.DeliveryAttempt
				}
				agg := msg.Attributes["aggregateId"]
				reqLog.Infof("Handling activity event: message_id=%s delivery_attempt=%d aggregate_id=%s publish_time=%s", msg.ID, attempt, agg, msg.PublishTime.Format(time.RFC3339))

				if err := dispatcher.Process(ctx, msg); err != nil {
					if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
						reqLog.Warnf("Handler context canceled: message_id=%s delivery_attempt=%d; nacking for redelivery", msg.ID, attempt)
					} else {
						reqLog.Errorf("Failed to handle activity event: error=%v, message_id=%s, request_id=%s", err, msg.ID, reqID)
					}
					msg.Nack()
				} else {
					reqLog.Infof("Successfully handled activity event: message_id=%s", msg.ID)
					msg.Ack()
				}
			})

			// Handle Receive termination conditions explicitly
			if err == nil {
				pubSubLog.Warnf("Pub/Sub Receive returned nil (no error); restarting receive loop")
				backoff = time.Second
			} else if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				pubSubLog.Infof("Pub/Sub Receive stopped due to context cancellation: %v; exiting subscriber loop", err)
				return
			} else {
				pubSubLog.Warnf("Pub/Sub Receive returned error: %v; will retry", err)
			}

			sleep := backoff
			if backoff < 30*time.Second {
				backoff *= 2
			}
			pubSubLog.Infof("Retrying subscriber after backoff=%s (attempt=%d)", sleep, attempt)
			select {
			case <-time.After(sleep):
				continue
			case <-ctxWithCancel.Done():
				pubSubLog.Infof("Context canceled during backoff; exiting subscriber loop")
				return
			}
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Infof("Activity Read Model Updater started successfully")
	<-sigChan

	log.Infof("Shutting down Activity Read Model Updater...")
	cancel()

	// Give some time for graceful shutdown
	time.Sleep(5 * time.Second)
	log.Infof("Activity Read Model Updater stopped")
}

func initFirestore(ctx context.Context, credentialsPath, projectID string) (*firestore.Client, error) {
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

func initPubSub(ctx context.Context, credentialsPath, projectID string) (*pubsub.Client, error) {
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
