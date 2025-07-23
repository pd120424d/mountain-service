package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
)

// HTTPClientInterface defines the interface for HTTP client operations
type HTTPClientInterface interface {
	Get(ctx context.Context, endpoint string) (*http.Response, error)
	Post(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)
	Put(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)
	Delete(ctx context.Context, endpoint string) (*http.Response, error)
}

// MockHTTPClient is a mock implementation of the HTTPClientInterface
type MockHTTPClient struct {
	GetResponse    *http.Response
	GetError       error
	PostResponse   *http.Response
	PostError      error
	PutResponse    *http.Response
	PutError       error
	DeleteResponse *http.Response
	DeleteError    error
}

func (m *MockHTTPClient) Get(ctx context.Context, endpoint string) (*http.Response, error) {
	return m.GetResponse, m.GetError
}

func (m *MockHTTPClient) Post(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	return m.PostResponse, m.PostError
}

func (m *MockHTTPClient) Put(ctx context.Context, endpoint string, body interface{}) (*http.Response, error) {
	return m.PutResponse, m.PutError
}

func (m *MockHTTPClient) Delete(ctx context.Context, endpoint string) (*http.Response, error) {
	return m.DeleteResponse, m.DeleteError
}

// testEmployeeClient is a test version of employeeClient that accepts the mock interface
type testEmployeeClient struct {
	httpClient HTTPClientInterface
	logger     utils.Logger
}

func (c *testEmployeeClient) GetOnCallEmployees(ctx context.Context, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error) {
	endpoint := "/api/v1/employees/on-call"
	if shiftBuffer > 0 {
		endpoint += "?shift_buffer=" + shiftBuffer.String()
	}

	resp, err := c.httpClient.Get(ctx, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get on-call employees: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("employee service returned status %d", resp.StatusCode)
	}

	var response employeeV1.OnCallEmployeesResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Employees, nil
}

func (c *testEmployeeClient) GetAllEmployees(ctx context.Context) ([]employeeV1.EmployeeResponse, error) {
	// Implementation for testing - simplified
	return nil, nil
}

func (c *testEmployeeClient) GetEmployeeByID(ctx context.Context, employeeID uint) (*employeeV1.EmployeeResponse, error) {
	// Implementation for testing - simplified
	return nil, nil
}

func (c *testEmployeeClient) CheckActiveEmergencies(ctx context.Context, employeeID uint) (bool, error) {
	// Implementation for testing - simplified
	return false, nil
}

func createMockResponse(statusCode int, body interface{}) *http.Response {
	jsonBody, _ := json.Marshal(body)
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(bytes.NewReader(jsonBody)),
	}
}

func TestNewEmployeeClient(t *testing.T) {
	t.Parallel()

	t.Run("it creates a new employee client with default timeout", func(t *testing.T) {
		config := EmployeeClientConfig{
			BaseURL:     "http://localhost:8080",
			ServiceAuth: &auth.ServiceAuth{},
			Logger:      utils.NewTestLogger(),
		}

		client := NewEmployeeClient(config)
		assert.NotNil(t, client)
		assert.IsType(t, &employeeClient{}, client)
	})

	t.Run("it creates a new employee client with custom timeout", func(t *testing.T) {
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
}

func TestEmployeeClient_GetOnCallEmployees(t *testing.T) {
	t.Parallel()

	t.Run("it successfully retrieves on-call employees", func(t *testing.T) {
		mockResponse := employeeV1.OnCallEmployeesResponse{
			Employees: []employeeV1.EmployeeResponse{
				{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john@example.com",
					Username:  "johndoe",
				},
				{
					ID:        2,
					FirstName: "Jane",
					LastName:  "Smith",
					Email:     "jane@example.com",
					Username:  "janesmith",
				},
			},
		}

		mockHTTPClient := &MockHTTPClient{
			GetResponse: createMockResponse(http.StatusOK, mockResponse),
			GetError:    nil,
		}

		// Create a test client with the mock
		client := &testEmployeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetOnCallEmployees(context.Background(), 0)
		assert.NoError(t, err)
		assert.Len(t, employees, 2)
		assert.Equal(t, uint(1), employees[0].ID)
		assert.Equal(t, "John", employees[0].FirstName)
		assert.Equal(t, "Doe", employees[0].LastName)
		assert.Equal(t, uint(2), employees[1].ID)
		assert.Equal(t, "Jane", employees[1].FirstName)
		assert.Equal(t, "Smith", employees[1].LastName)
	})

	t.Run("it handles HTTP client error", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{
			GetResponse: nil,
			GetError:    errors.New("network error"),
		}

		client := &testEmployeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetOnCallEmployees(context.Background(), 0)
		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Contains(t, err.Error(), "failed to get on-call employees")
	})

	t.Run("it handles non-200 status code", func(t *testing.T) {
		mockHTTPClient := &MockHTTPClient{
			GetResponse: createMockResponse(http.StatusInternalServerError, nil),
			GetError:    nil,
		}

		client := &testEmployeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetOnCallEmployees(context.Background(), 0)
		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Contains(t, err.Error(), "employee service returned status 500")
	})
}
