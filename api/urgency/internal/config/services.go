package config

import (
	"os"
	"time"

	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/pd120424d/mountain-service/api/urgency/internal/clients"
)

// ServiceClients holds all external service clients
type ServiceClients struct {
	EmployeeClient clients.EmployeeClient
	// TODO: Add activity client here as well
}

// ServiceConfig holds configuration for external services
type ServiceConfig struct {
	EmployeeServiceURL string
	ActivityServiceURL string
	ServiceAuthSecret  string
	ServiceName        string
}

// LoadServiceConfig loads service configuration from environment variables
func LoadServiceConfig() ServiceConfig {
	return ServiceConfig{
		EmployeeServiceURL: getEnvOrDefault("EMPLOYEE_SERVICE_URL", "http://employee-service:8082"),
		ActivityServiceURL: getEnvOrDefault("ACTIVITY_SERVICE_URL", "http://activity-service:8084"),
		ServiceAuthSecret:  getEnvOrDefault("SERVICE_AUTH_SECRET", "super-secret-service-auth-key"),
		ServiceName:        "urgency-service",
	}
}

// InitializeServiceClients creates and configures all external service clients
func InitializeServiceClients(config ServiceConfig, logger utils.Logger) (*ServiceClients, error) {
	// Create service authentication
	serviceAuth := auth.NewServiceAuth(auth.ServiceAuthConfig{
		Secret:      config.ServiceAuthSecret,
		ServiceName: config.ServiceName,
		TokenTTL:    time.Hour,
	})

	// Initialize employee client
	employeeClient := clients.NewEmployeeClient(clients.EmployeeClientConfig{
		BaseURL:     config.EmployeeServiceURL,
		ServiceAuth: serviceAuth,
		Logger:      logger,
		Timeout:     30 * time.Second,
	})

	return &ServiceClients{
		EmployeeClient: employeeClient,
	}, nil
}

// getEnvOrDefault returns the environment variable value or a default value if not set
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
