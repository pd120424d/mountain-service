package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockHTTPClient implements the HTTPClient interface for testing
type MockHTTPClient struct {
	GetFunc    func(ctx context.Context, endpoint string) (*http.Response, error)
	PostFunc   func(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)
	PutFunc    func(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)
	DeleteFunc func(ctx context.Context, endpoint string) (*http.Response, error)
	
	// Track calls for verification
	GetCalls    []string
	PostCalls   []string
	PutCalls    []string
	DeleteCalls []string
}

func (m *MockHTTPClient) Get(ctx context.Context, endpoint string) (*http.Response, error) {
	m.GetCalls = append(m.GetCalls, endpoint)
	if m.GetFunc != nil {
		return m.GetFunc(ctx, endpoint)
	}
	return nil, fmt.Errorf("mock not configured for GET %s", endpoint)
}

func (m *MockHTTPClient) Post(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	m.PostCalls = append(m.PostCalls, endpoint)
	if m.PostFunc != nil {
		return m.PostFunc(ctx, endpoint, body)
	}
	return nil, fmt.Errorf("mock not configured for POST %s", endpoint)
}

func (m *MockHTTPClient) Put(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	m.PutCalls = append(m.PutCalls, endpoint)
	if m.PutFunc != nil {
		return m.PutFunc(ctx, endpoint, body)
	}
	return nil, fmt.Errorf("mock not configured for PUT %s", endpoint)
}

func (m *MockHTTPClient) Delete(ctx context.Context, endpoint string) (*http.Response, error) {
	m.DeleteCalls = append(m.DeleteCalls, endpoint)
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, endpoint)
	}
	return nil, fmt.Errorf("mock not configured for DELETE %s", endpoint)
}

// Helper function to create a mock HTTP response
func createMockResponse(statusCode int, body interface{}) *http.Response {
	jsonBody, _ := json.Marshal(body)
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(jsonBody)),
	}
}

// Helper function to create an employee client with mock HTTP client for testing
func createTestEmployeeClient(mockHTTPClient *MockHTTPClient) EmployeeClient {
	return NewEmployeeClientWithHTTPClient(mockHTTPClient, utils.NewTestLogger())
}

func TestNewEmployeeClient(t *testing.T) {
	t.Parallel()

	t.Run("creates a new employee client with default timeout", func(t *testing.T) {
		config := EmployeeClientConfig{
			BaseURL:     "http://localhost:8080",
			ServiceAuth: &auth.ServiceAuth{},
			Logger:      utils.NewTestLogger(),
		}

		client := NewEmployeeClient(config)
		assert.NotNil(t, client)
		assert.IsType(t, &employeeClient{}, client)
	})

	t.Run("creates a new employee client with custom timeout", func(t *testing.T) {
		config := EmployeeClientConfig{
			BaseURL:     "http://localhost:8080",
			ServiceAuth: &auth.ServiceAuth{},
			Logger:      utils.NewTestLogger(),
			Timeout:     60 * time.Second,
		}

		client := NewEmployeeClient(config)
		assert.NotNil(t, client)
		assert.IsType(t, &employeeClient{}, client)
	})

	t.Run("creates a new employee client with zero timeout defaults to 30 seconds", func(t *testing.T) {
		config := EmployeeClientConfig{
			BaseURL:     "http://localhost:8080",
			ServiceAuth: &auth.ServiceAuth{},
			Logger:      utils.NewTestLogger(),
			Timeout:     0, // Should default to 30 seconds
		}

		client := NewEmployeeClient(config)
		assert.NotNil(t, client)
		assert.IsType(t, &employeeClient{}, client)
	})
}

func TestNewEmployeeClientWithHTTPClient(t *testing.T) {
	t.Parallel()

	t.Run("creates a new employee client with injected HTTP client", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{}
		logger := utils.NewTestLogger()

		client := NewEmployeeClientWithHTTPClient(mockHTTPClient, logger)
		
		assert.NotNil(t, client)
		assert.IsType(t, &employeeClient{}, client)
	})
}

func TestEmployeeClient_GetOnCallEmployees(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves on-call employees without shift buffer", func(t *testing.T) {
		mockResponse := employeeV1.OnCallEmployeesResponse{
			Employees: []employeeV1.EmployeeResponse{
				{
					ID:          1,
					FirstName:   "John",
					LastName:    "Doe",
					Email:       "john@example.com",
					Username:    "johndoe",
					Gender:      "Male",
					Phone:       "+1234567890",
					ProfileType: "Medic",
				},
				{
					ID:          2,
					FirstName:   "Jane",
					LastName:    "Smith",
					Email:       "jane@example.com",
					Username:    "janesmith",
					Gender:      "Female",
					Phone:       "+1234567891",
					ProfileType: "Technical",
				},
			},
		}

		mockHTTPClient := &MockHTTPClient{
			GetFunc: func(ctx context.Context, endpoint string) (*http.Response, error) {
				assert.Equal(t, "/api/v1/employees/on-call", endpoint)
				return createMockResponse(http.StatusOK, mockResponse), nil
			},
		}

		client := createTestEmployeeClient(mockHTTPClient)
		employees, err := client.GetOnCallEmployees(context.Background(), 0)
		
		require.NoError(t, err)
		assert.Len(t, employees, 2)
		assert.Equal(t, uint(1), employees[0].ID)
		assert.Equal(t, "John", employees[0].FirstName)
		assert.Equal(t, "Doe", employees[0].LastName)
		assert.Equal(t, uint(2), employees[1].ID)
		assert.Equal(t, "Jane", employees[1].FirstName)
		assert.Equal(t, "Smith", employees[1].LastName)
		assert.Len(t, mockHTTPClient.GetCalls, 1)
	})

	t.Run("successfully retrieves on-call employees with shift buffer", func(t *testing.T) {
		mockResponse := employeeV1.OnCallEmployeesResponse{
			Employees: []employeeV1.EmployeeResponse{
				{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@example.com",
					Username:  "johndoe",
				},
			},
		}

		mockHTTPClient := &MockHTTPClient{
			GetFunc: func(ctx context.Context, endpoint string) (*http.Response, error) {
				assert.Equal(t, "/api/v1/employees/on-call?shift_buffer=1h0m0s", endpoint)
				return createMockResponse(http.StatusOK, mockResponse), nil
			},
		}

		client := createTestEmployeeClient(mockHTTPClient)
		employees, err := client.GetOnCallEmployees(context.Background(), time.Hour)
		
		require.NoError(t, err)
		assert.Len(t, employees, 1)
		assert.Equal(t, uint(1), employees[0].ID)
		assert.Len(t, mockHTTPClient.GetCalls, 1)
	})

	t.Run("handles HTTP client error", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{
			GetFunc: func(ctx context.Context, endpoint string) (*http.Response, error) {
				return nil, fmt.Errorf("network error")
			},
		}

		client := createTestEmployeeClient(mockHTTPClient)
		employees, err := client.GetOnCallEmployees(context.Background(), 0)
		
		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Contains(t, err.Error(), "failed to get on-call employees")
		assert.Contains(t, err.Error(), "network error")
	})

	t.Run("handles non-200 status code", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{
			GetFunc: func(ctx context.Context, endpoint string) (*http.Response, error) {
				return createMockResponse(http.StatusInternalServerError, nil), nil
			},
		}

		client := createTestEmployeeClient(mockHTTPClient)
		employees, err := client.GetOnCallEmployees(context.Background(), 0)
		
		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Contains(t, err.Error(), "employee service returned status 500")
	})

	t.Run("handles JSON decode error", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{
			GetFunc: func(ctx context.Context, endpoint string) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
				}, nil
			},
		}

		client := createTestEmployeeClient(mockHTTPClient)
		employees, err := client.GetOnCallEmployees(context.Background(), 0)
		
		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Contains(t, err.Error(), "failed to decode response")
	})
}
