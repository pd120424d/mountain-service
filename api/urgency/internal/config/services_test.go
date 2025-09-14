package config

import (
	"testing"

	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

func TestLoadServiceConfig(t *testing.T) {
	// t.Parallel() disabled: subtests use t.Setenv

	t.Run("it returns default values when environment variables are not set", func(t *testing.T) {
		config := LoadServiceConfig()
		assert.Equal(t, "http://employee-service:8082", config.EmployeeServiceURL)
		assert.Equal(t, "http://activity-service:8084", config.ActivityServiceURL)
		assert.Equal(t, "super-secret-service-auth-key", config.ServiceAuthSecret)
		assert.Equal(t, "urgency-service", config.ServiceName)
	})

	t.Run("it returns environment variable values when set", func(t *testing.T) {
		t.Setenv("EMPLOYEE_SERVICE_URL", "http://test-employee-service")
		t.Setenv("ACTIVITY_SERVICE_URL", "http://test-activity-service")
		t.Setenv("SERVICE_AUTH_SECRET", "test-secret")
		config := LoadServiceConfig()
		assert.Equal(t, "http://test-employee-service", config.EmployeeServiceURL)
		assert.Equal(t, "http://test-activity-service", config.ActivityServiceURL)
		assert.Equal(t, "test-secret", config.ServiceAuthSecret)
		assert.Equal(t, "urgency-service", config.ServiceName)
	})
}

func TestInitializeServiceClients(t *testing.T) {
	t.Run("it initializes service clients successfully", func(t *testing.T) {
		config := LoadServiceConfig()
		logger := utils.NewTestLogger()
		clients, err := InitializeServiceClients(config, logger)
		assert.NoError(t, err)
		assert.NotNil(t, clients.EmployeeClient)
	})
}
