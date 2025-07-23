package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/client"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"go.uber.org/zap"
)

// HTTPClient interface defines the HTTP operations needed by the employee client
type HTTPClient interface {
	Get(ctx context.Context, endpoint string) (*http.Response, error)
	Post(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)
	Put(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)
	Delete(ctx context.Context, endpoint string) (*http.Response, error)
}

// EmployeeClient interface for communicating with the employee service
type EmployeeClient interface {
	GetOnCallEmployees(ctx context.Context, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error)
	GetAllEmployees(ctx context.Context) ([]employeeV1.EmployeeResponse, error)
	GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error)
	CheckActiveEmergencies(ctx context.Context, employeeID uint) (bool, error)
}

// employeeClient implements the EmployeeClient interface
type employeeClient struct {
	httpClient HTTPClient
	logger     utils.Logger
}

// EmployeeClientConfig holds the configuration for the employee client
type EmployeeClientConfig struct {
	BaseURL     string
	ServiceAuth *auth.ServiceAuth
	Logger      utils.Logger
	Timeout     time.Duration
}

// NewEmployeeClient creates a new employee service client
func NewEmployeeClient(config EmployeeClientConfig) EmployeeClient {
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	httpClient := client.NewHTTPClient(client.HTTPClientConfig{
		BaseURL:     config.BaseURL,
		Timeout:     config.Timeout,
		ServiceAuth: config.ServiceAuth,
		Logger:      config.Logger,
	})

	return NewEmployeeClientWithHTTPClient(httpClient, config.Logger)
}

// NewEmployeeClientWithHTTPClient creates a new employee service client with a custom HTTP client
// This constructor is useful for testing as it allows dependency injection
func NewEmployeeClientWithHTTPClient(httpClient HTTPClient, logger utils.Logger) EmployeeClient {
	return &employeeClient{
		httpClient: httpClient,
		logger:     logger.WithName("employeeClient"),
	}
}

// GetOnCallEmployees retrieves employees currently on call
func (c *employeeClient) GetOnCallEmployees(ctx context.Context, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error) {
	c.logger.Info("Getting on-call employees", zap.Duration("shiftBuffer", shiftBuffer))

	endpoint := "/api/v1/employees/on-call"
	if shiftBuffer > 0 {
		endpoint += fmt.Sprintf("?shift_buffer=%s", shiftBuffer.String())
	}

	resp, err := c.httpClient.Get(ctx, endpoint)
	if err != nil {
		c.logger.Error("Failed to get on-call employees", zap.Error(err))
		return nil, fmt.Errorf("failed to get on-call employees: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Employee service returned error", zap.Int("statusCode", resp.StatusCode))
		return nil, fmt.Errorf("employee service returned status %d", resp.StatusCode)
	}

	var response employeeV1.OnCallEmployeesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("Failed to decode on-call employees response", zap.Error(err))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Info("Successfully retrieved on-call employees", zap.Int("count", len(response.Employees)))
	return response.Employees, nil
}

// GetAllEmployees retrieves all employees
func (c *employeeClient) GetAllEmployees(ctx context.Context) ([]employeeV1.EmployeeResponse, error) {
	c.logger.Info("Getting all employees")

	resp, err := c.httpClient.Get(ctx, "/api/v1/employees")
	if err != nil {
		c.logger.Error("Failed to get all employees", zap.Error(err))
		return nil, fmt.Errorf("failed to get all employees: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Employee service returned error", zap.Int("statusCode", resp.StatusCode))
		return nil, fmt.Errorf("employee service returned status %d", resp.StatusCode)
	}

	var response employeeV1.AllEmployeesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("Failed to decode all employees response", zap.Error(err))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Info("Successfully retrieved all employees", zap.Int("count", len(response.Employees)))
	return response.Employees, nil
}

// GetEmployeeByID retrieves a specific employee by ID
func (c *employeeClient) GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error) {
	c.logger.Info("Getting employee by ID", zap.Uint("employeeID", employeeID))

	endpoint := fmt.Sprintf("/api/v1/employees/%d", employeeID)
	resp, err := c.httpClient.Get(ctx, endpoint)
	if err != nil {
		c.logger.Error("Failed to get employee by ID", zap.Uint("employeeID", employeeID), zap.Error(err))
		return nil, fmt.Errorf("failed to get employee %d: %w", employeeID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		c.logger.Warn("Employee not found", zap.Uint("employeeID", employeeID))
		return nil, fmt.Errorf("employee %d not found", employeeID)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Employee service returned error", zap.Int("statusCode", resp.StatusCode), zap.Uint("employeeID", employeeID))
		return nil, fmt.Errorf("employee service returned status %d", resp.StatusCode)
	}

	var employee employeeV1.EmployeeResponse
	if err := json.NewDecoder(resp.Body).Decode(&employee); err != nil {
		c.logger.Error("Failed to decode employee response", zap.Uint("employeeID", employeeID), zap.Error(err))
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Info("Successfully retrieved employee", zap.Uint("employeeID", employeeID))
	return &employee, nil
}

// CheckActiveEmergencies checks if an employee has active emergencies
func (c *employeeClient) CheckActiveEmergencies(ctx context.Context, employeeID uint) (bool, error) {
	c.logger.Info("Checking active emergencies for employee", zap.Uint("employeeID", employeeID))

	endpoint := fmt.Sprintf("/api/v1/employees/%d/active-emergencies", employeeID)
	resp, err := c.httpClient.Get(ctx, endpoint)
	if err != nil {
		c.logger.Error("Failed to check active emergencies", zap.Uint("employeeID", employeeID), zap.Error(err))
		return false, fmt.Errorf("failed to check active emergencies for employee %d: %w", employeeID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Employee service returned error", zap.Int("statusCode", resp.StatusCode), zap.Uint("employeeID", employeeID))
		return false, fmt.Errorf("employee service returned status %d", resp.StatusCode)
	}

	var response employeeV1.ActiveEmergenciesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		c.logger.Error("Failed to decode active emergencies response", zap.Uint("employeeID", employeeID), zap.Error(err))
		return false, fmt.Errorf("failed to decode response: %w", err)
	}

	c.logger.Info("Successfully checked active emergencies", zap.Uint("employeeID", employeeID), zap.Bool("hasActive", response.HasActiveEmergencies))
	return response.HasActiveEmergencies, nil
}
