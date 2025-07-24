package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

const (
	// ServerPort is the port the activity service runs on
	ServerPort = 8083

	// DatabaseName is the name of the activity database
	DatabaseName = "activity_service"

	// DefaultPageSize is the default number of items per page
	DefaultPageSize = 50

	// MaxPageSize is the maximum number of items per page
	MaxPageSize = 1000
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

// GetDefaultDatabaseConfig returns default database configuration
func GetDefaultDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASSWORD", "postgres"),
		DBName:   getEnvOrDefault("DB_NAME", "mountain_service"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		TimeZone: getEnvOrDefault("DB_TIMEZONE", "UTC"),
	}
}

// GetConnectionString returns the database connection string
func (c *DatabaseConfig) GetConnectionString() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		c.Host, c.User, c.Password, c.DBName, c.Port, c.SSLMode, c.TimeZone,
	)
}

// GetActivityDB initializes and returns a database connection for the activity service
func GetActivityDB(logger utils.Logger, connectionString string) *gorm.DB {
	logger.Info("Connecting to activity database")

	// Configure GORM logger
	dbLogger := gormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  gormLogger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)

	// Open database connection
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: dbLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
	})

	if err != nil {
		logger.Fatalf("Failed to connect to activity database: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		logger.Fatalf("Failed to get underlying sql.DB: %v", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Test the connection
	if err := sqlDB.Ping(); err != nil {
		logger.Fatalf("Failed to ping activity database: %v", err)
	}

	logger.Info("Successfully connected to activity database")
	return db
}

// ServiceConfig holds configuration for the activity service
type ServiceConfig struct {
	Port            int
	DatabaseConfig  *DatabaseConfig
	JWTSecret       string
	LogLevel        string
	Environment     string
	ShutdownTimeout time.Duration
}

// GetDefaultServiceConfig returns default service configuration
func GetDefaultServiceConfig() *ServiceConfig {
	return &ServiceConfig{
		Port:            ServerPort,
		DatabaseConfig:  GetDefaultDatabaseConfig(),
		JWTSecret:       getEnvOrDefault("JWT_SECRET", ""),
		LogLevel:        getEnvOrDefault("LOG_LEVEL", "info"),
		Environment:     getEnvOrDefault("ENVIRONMENT", "development"),
		ShutdownTimeout: 5 * time.Second,
	}
}

// Validate validates the service configuration
func (c *ServiceConfig) Validate() error {
	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("invalid port: %d", c.Port)
	}

	if c.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	if c.DatabaseConfig == nil {
		return fmt.Errorf("database configuration is required")
	}

	return nil
}

// getEnvOrDefault returns the value of an environment variable or a default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// LoadConfig loads configuration from environment variables
func LoadConfig() (*ServiceConfig, error) {
	config := GetDefaultServiceConfig()

	// Override with environment variables if present
	if port := os.Getenv("PORT"); port != "" {
		var portInt int
		if _, err := fmt.Sscanf(port, "%d", &portInt); err == nil {
			config.Port = portInt
		}
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("configuration validation failed: %w", err)
	}

	return config, nil
}

// LogConfig logs the current configuration (without sensitive data)
func LogConfig(config *ServiceConfig) {
	log.Printf("Activity Service Configuration:")
	log.Printf("  Port: %d", config.Port)
	log.Printf("  Database Host: %s", config.DatabaseConfig.Host)
	log.Printf("  Database Port: %s", config.DatabaseConfig.Port)
	log.Printf("  Database Name: %s", config.DatabaseConfig.DBName)
	log.Printf("  Log Level: %s", config.LogLevel)
	log.Printf("  Environment: %s", config.Environment)
	log.Printf("  Shutdown Timeout: %v", config.ShutdownTimeout)
}
