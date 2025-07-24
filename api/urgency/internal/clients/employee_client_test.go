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
	"go.uber.org/mock/gomock"
)

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

	t.Run("it creates a new employee client with zero timeout defaults to 30 seconds", func(t *testing.T) {
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

	t.Run("it creates a new employee client with injected HTTP client", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := NewEmployeeClientWithHTTPClient(mockHTTPClient, logger.WithName("employeeClient"))

		assert.NotNil(t, client)
		assert.IsType(t, &employeeClient{}, client)
	})
}

func TestEmployeeClient_GetOnCallEmployees(t *testing.T) {
	t.Parallel()

	t.Run("it successfully retrieves on-call employees without shift buffer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

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

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/on-call").Return(createMockResponse(http.StatusOK, mockResponse), nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetOnCallEmployees(context.Background(), 0)

		require.NoError(t, err)
		assert.Len(t, employees, 2)
		assert.Equal(t, uint(1), employees[0].ID)
		assert.Equal(t, "John", employees[0].FirstName)
		assert.Equal(t, "Doe", employees[0].LastName)
		assert.Equal(t, uint(2), employees[1].ID)
		assert.Equal(t, "Jane", employees[1].FirstName)
		assert.Equal(t, "Smith", employees[1].LastName)
	})

	t.Run("it successfully retrieves on-call employees with shift buffer", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

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

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/on-call?shift_buffer=1h0m0s").Return(createMockResponse(http.StatusOK, mockResponse), nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetOnCallEmployees(context.Background(), time.Hour)

		require.NoError(t, err)
		assert.Len(t, employees, 1)
		assert.Equal(t, uint(1), employees[0].ID)
	})

	t.Run("it handles HTTP client error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/on-call").Return(nil, fmt.Errorf("network error"))

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetOnCallEmployees(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Contains(t, err.Error(), "failed to get on-call employees")
		assert.Contains(t, err.Error(), "network error")
	})

	t.Run("it handles non-200 status code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/on-call").Return(createMockResponse(http.StatusInternalServerError, nil), nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetOnCallEmployees(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Contains(t, err.Error(), "employee service returned status 500")
	})

	t.Run("it handles JSON decode error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/on-call").Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetOnCallEmployees(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Contains(t, err.Error(), "failed to decode response")
	})
}

func TestEmployeeClient_GetAllEmployees(t *testing.T) {
	t.Parallel()

	t.Run("it successfully retrieves all employees", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockResponse := employeeV1.AllEmployeesResponse{
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

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees").Return(createMockResponse(http.StatusOK, mockResponse), nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetAllEmployees(context.Background())

		require.NoError(t, err)
		assert.Len(t, employees, 2)
		assert.Equal(t, uint(1), employees[0].ID)
		assert.Equal(t, "John", employees[0].FirstName)
		assert.Equal(t, "Doe", employees[0].LastName)
		assert.Equal(t, uint(2), employees[1].ID)
		assert.Equal(t, "Jane", employees[1].FirstName)
		assert.Equal(t, "Smith", employees[1].LastName)
	})

	t.Run("handles HTTP client error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees").Return(nil, fmt.Errorf("network error"))

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetAllEmployees(context.Background())

		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Contains(t, err.Error(), "failed to get all employees")
		assert.Contains(t, err.Error(), "network error")
	})

	t.Run("handles non-200 status code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees").Return(createMockResponse(http.StatusInternalServerError, nil), nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetAllEmployees(context.Background())

		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Contains(t, err.Error(), "employee service returned status 500")
	})

	t.Run("handles JSON decode error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees").Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employees, err := client.GetAllEmployees(context.Background())

		assert.Error(t, err)
		assert.Nil(t, employees)
		assert.Contains(t, err.Error(), "failed to decode response")
	})
}

func TestEmployeeClient_GetEmployeeByID(t *testing.T) {
	t.Parallel()

	t.Run("it successfully retrieves employee by ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockResponse := employeeV1.EmployeeResponse{
			ID:          1,
			FirstName:   "John",
			LastName:    "Doe",
			Email:       "john@example.com",
			Username:    "johndoe",
			Gender:      "Male",
			Phone:       "+1234567890",
			ProfileType: "Medic",
		}

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/1").Return(createMockResponse(http.StatusOK, mockResponse), nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employee, err := client.GetEmployeeByID(context.Background(), 1)

		require.NoError(t, err)
		require.NotNil(t, employee)
		assert.Equal(t, uint(1), employee.ID)
		assert.Equal(t, "John", employee.FirstName)
		assert.Equal(t, "Doe", employee.LastName)
		assert.Equal(t, "john@example.com", employee.Email)
	})

	t.Run("it handles HTTP client error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/1").Return(nil, fmt.Errorf("network error"))

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employee, err := client.GetEmployeeByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, employee)
		assert.Contains(t, err.Error(), "failed to get employee")
		assert.Contains(t, err.Error(), "network error")
	})

	t.Run("it handles non-200 status code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/999").Return(createMockResponse(http.StatusInternalServerError, nil), nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employee, err := client.GetEmployeeByID(context.Background(), 999)

		assert.Error(t, err)
		assert.Nil(t, employee)
		assert.Contains(t, err.Error(), "employee service returned status 500")
	})

	t.Run("it handles 404 not found", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/999").Return(createMockResponse(http.StatusNotFound, nil), nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employee, err := client.GetEmployeeByID(context.Background(), 999)

		assert.Error(t, err)
		assert.Nil(t, employee)
		assert.Contains(t, err.Error(), "employee 999 not found")
	})

	t.Run("it handles JSON decode error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/1").Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		employee, err := client.GetEmployeeByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, employee)
		assert.Contains(t, err.Error(), "failed to decode response")
	})
}

func TestEmployeeClient_CheckActiveEmergencies(t *testing.T) {
	t.Parallel()

	t.Run("it successfully checks active emergencies - has active", func(t *testing.T) {
		mockResponse := employeeV1.ActiveEmergenciesResponse{
			HasActiveEmergencies: true,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/1/active-emergencies").Return(createMockResponse(http.StatusOK, mockResponse), nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		hasActive, err := client.CheckActiveEmergencies(context.Background(), 1)

		require.NoError(t, err)
		assert.True(t, hasActive)
	})

	t.Run("it successfully checks active emergencies - no active", func(t *testing.T) {
		mockResponse := employeeV1.ActiveEmergenciesResponse{
			HasActiveEmergencies: false,
		}

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/2/active-emergencies").Return(createMockResponse(http.StatusOK, mockResponse), nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		hasActive, err := client.CheckActiveEmergencies(context.Background(), 2)

		require.NoError(t, err)
		assert.False(t, hasActive)
	})

	t.Run("it handles HTTP client error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/1/active-emergencies").Return(nil, fmt.Errorf("network error"))

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		hasActive, err := client.CheckActiveEmergencies(context.Background(), 1)

		assert.Error(t, err)
		assert.False(t, hasActive)
		assert.Contains(t, err.Error(), "failed to check active emergencies")
		assert.Contains(t, err.Error(), "network error")
	})

	t.Run("it handles non-200 status code", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/1/active-emergencies").Return(createMockResponse(http.StatusInternalServerError, nil), nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		hasActive, err := client.CheckActiveEmergencies(context.Background(), 1)

		assert.Error(t, err)
		assert.False(t, hasActive)
		assert.Contains(t, err.Error(), "employee service returned status 500")
	})

	t.Run("it handles JSON decode error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTPClient := NewMockHTTPClient(ctrl)
		mockHTTPClient.EXPECT().Get(gomock.Any(), "/api/v1/employees/1/active-emergencies").Return(&http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte("invalid json"))),
		}, nil)

		client := &employeeClient{
			httpClient: mockHTTPClient,
			logger:     utils.NewTestLogger(),
		}

		hasActive, err := client.CheckActiveEmergencies(context.Background(), 1)

		assert.Error(t, err)
		assert.False(t, hasActive)
		assert.Contains(t, err.Error(), "failed to decode response")
	})
}

func createMockResponse(statusCode int, body interface{}) *http.Response {
	var bodyReader io.ReadCloser
	if body != nil {
		jsonBytes, _ := json.Marshal(body)
		bodyReader = io.NopCloser(bytes.NewReader(jsonBytes))
	} else {
		bodyReader = io.NopCloser(bytes.NewReader([]byte{}))
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       bodyReader,
	}
}
