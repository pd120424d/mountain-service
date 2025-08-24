package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig_LoadConfig(t *testing.T) {
	t.Parallel()

	t.Run("it succeeds when loading default config", func(t *testing.T) {
		config := loadConfig()
		assert.NotNil(t, config)
		assert.Equal(t, "postgres://user:password@localhost:5432/activities?sslmode=disable", config.DatabaseURL)
		assert.Equal(t, "your-project-id", config.FirebaseProjectID)
		assert.Equal(t, "activity-events", config.PubSubTopic)
		assert.Equal(t, "activity-events-sub", config.PubSubSubscription)
		assert.Equal(t, 10, config.OutboxPollIntervalSeconds)
		assert.Equal(t, 8090, config.HealthPort)
		assert.Equal(t, "info", config.LogLevel)
		assert.Equal(t, "dev", config.Version)
		assert.Equal(t, "unknown", config.GitSHA)
	})
}

func TestGetEnvOrDefault(t *testing.T) {
	t.Parallel()

	t.Run("it returns default value when env var is not set", func(t *testing.T) {
		result := getEnvOrDefault("NON_EXISTENT_VAR", "default")
		assert.Equal(t, "default", result)
	})
}

func TestGetEnvAsIntOrDefault(t *testing.T) {
	t.Parallel()

	t.Run("it returns default value when env var is not set", func(t *testing.T) {
		result := getEnvAsIntOrDefault("NON_EXISTENT_VAR", 42)
		assert.Equal(t, 42, result)
	})
}

func TestMain_Integration(t *testing.T) {
	t.Parallel()

	t.Run("it loads configuration correctly", func(t *testing.T) {
		// Test that main configuration loading works
		config := loadConfig()
		assert.NotNil(t, config)
		assert.NotEmpty(t, config.DatabaseURL)
		assert.NotEmpty(t, config.FirebaseProjectID)
		assert.NotEmpty(t, config.PubSubTopic)
		assert.NotEmpty(t, config.PubSubSubscription)
		assert.Greater(t, config.OutboxPollIntervalSeconds, 0)
		assert.Greater(t, config.HealthPort, 0)
		assert.NotEmpty(t, config.LogLevel)
		assert.NotEmpty(t, config.Version)
		assert.NotEmpty(t, config.GitSHA)
	})
}

func TestConfig_Validation(t *testing.T) {
	t.Parallel()

	t.Run("it has reasonable default values", func(t *testing.T) {
		config := loadConfig()

		// Validate default values are reasonable
		assert.Contains(t, config.DatabaseURL, "postgres://")
		assert.Equal(t, "activity-events", config.PubSubTopic)
		assert.Equal(t, "activity-events-sub", config.PubSubSubscription)
		assert.Equal(t, 10, config.OutboxPollIntervalSeconds)
		assert.Equal(t, 8090, config.HealthPort)
		assert.Equal(t, "info", config.LogLevel)
		assert.Equal(t, "dev", config.Version)
		assert.Equal(t, "unknown", config.GitSHA)
	})
}

func TestUpdater_EnvironmentHandling(t *testing.T) {
	t.Parallel()

	t.Run("it handles environment variables correctly", func(t *testing.T) {
		// Test getEnvOrDefault function
		result := getEnvOrDefault("NON_EXISTENT_VAR", "default_value")
		assert.Equal(t, "default_value", result)
	})

	t.Run("it handles integer environment variables correctly", func(t *testing.T) {
		// Test getEnvAsIntOrDefault function
		result := getEnvAsIntOrDefault("NON_EXISTENT_INT_VAR", 42)
		assert.Equal(t, 42, result)
	})

	t.Run("it validates configuration values", func(t *testing.T) {
		config := loadConfig()

		// Validate that all required configuration fields are set
		assert.NotEmpty(t, config.DatabaseURL)
		assert.NotEmpty(t, config.FirebaseProjectID)
		assert.NotEmpty(t, config.PubSubTopic)
		assert.NotEmpty(t, config.PubSubSubscription)
		assert.Greater(t, config.OutboxPollIntervalSeconds, 0)
		assert.Greater(t, config.HealthPort, 0)
		assert.NotEmpty(t, config.LogLevel)
		assert.NotEmpty(t, config.Version)
		assert.NotEmpty(t, config.GitSHA)
	})
}
