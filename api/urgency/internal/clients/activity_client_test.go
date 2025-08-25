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

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestActivityClient_CreateActivity(t *testing.T) {
	t.Parallel()

	t.Run("successfully creates activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		req := &activityV1.ActivityCreateRequest{
			Description: "Test activity",
			EmployeeID:  123,
			UrgencyID:   456,
		}

		expectedResponse := &activityV1.ActivityResponse{
			ID:          1,
			Description: "Test activity",
			EmployeeID:  123,
			UrgencyID:   456,
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		}

		responseBody, _ := json.Marshal(expectedResponse)
		mockResp := &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}

		mockHTTP.EXPECT().Post(gomock.Any(), "/activities", req).Return(mockResp, nil)

		result, err := client.CreateActivity(context.Background(), req)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "Test activity", result.Description)
		assert.Equal(t, uint(123), result.EmployeeID)
		assert.Equal(t, uint(456), result.UrgencyID)
	})

	t.Run("returns error when HTTP request fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		req := &activityV1.ActivityCreateRequest{
			Description: "Test activity",
			EmployeeID:  123,
			UrgencyID:   456,
		}

		mockHTTP.EXPECT().Post(gomock.Any(), "/activities", req).Return(nil, fmt.Errorf("network error"))

		result, err := client.CreateActivity(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to create activity")
	})

	t.Run("returns error when service returns non-201 status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		req := &activityV1.ActivityCreateRequest{
			Description: "Test activity",
			EmployeeID:  123,
			UrgencyID:   456,
		}

		mockResp := &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"error":"validation failed"}`))),
		}

		mockHTTP.EXPECT().Post(gomock.Any(), "/activities", req).Return(mockResp, nil)

		result, err := client.CreateActivity(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "activity service returned status 400")
	})

	t.Run("returns error when response JSON is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		req := &activityV1.ActivityCreateRequest{
			Description: "Test activity",
			EmployeeID:  123,
			UrgencyID:   456,
		}

		mockResp := &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader([]byte(`invalid json`))),
		}

		mockHTTP.EXPECT().Post(gomock.Any(), "/activities", req).Return(mockResp, nil)

		result, err := client.CreateActivity(context.Background(), req)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to decode activity response")
	})
}

func TestActivityClient_GetActivitiesByUrgency(t *testing.T) {
	t.Parallel()

	t.Run("successfully retrieves activities", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		expectedActivities := []activityV1.ActivityResponse{
			{
				ID:          1,
				Description: "First activity",
				EmployeeID:  123,
				UrgencyID:   456,
				CreatedAt:   time.Now().Format(time.RFC3339),
				UpdatedAt:   time.Now().Format(time.RFC3339),
			},
			{
				ID:          2,
				Description: "Second activity",
				EmployeeID:  124,
				UrgencyID:   456,
				CreatedAt:   time.Now().Format(time.RFC3339),
				UpdatedAt:   time.Now().Format(time.RFC3339),
			},
		}

		listResponse := activityV1.ActivityListResponse{
			Activities: expectedActivities,
			Total:      2,
			Page:       1,
			PageSize:   10,
			TotalPages: 1,
		}

		responseBody, _ := json.Marshal(listResponse)
		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}

		mockHTTP.EXPECT().Get(gomock.Any(), "/activities?urgency_id=456").Return(mockResp, nil)

		result, err := client.GetActivitiesByUrgency(context.Background(), 456)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result, 2)
		assert.Equal(t, uint(1), result[0].ID)
		assert.Equal(t, "First activity", result[0].Description)
		assert.Equal(t, uint(456), result[0].UrgencyID)
	})

	t.Run("returns error when HTTP request fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		mockHTTP.EXPECT().Get(gomock.Any(), "/activities?urgency_id=456").Return(nil, fmt.Errorf("network error"))

		result, err := client.GetActivitiesByUrgency(context.Background(), 456)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to get activities")
	})

	t.Run("returns error when service returns non-200 status", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		mockResp := &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"error":"database error"}`))),
		}

		mockHTTP.EXPECT().Get(gomock.Any(), "/activities?urgency_id=456").Return(mockResp, nil)

		result, err := client.GetActivitiesByUrgency(context.Background(), 456)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "activity service returned status 500")
	})

	t.Run("returns error when response JSON is invalid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(`invalid json`))),
		}

		mockHTTP.EXPECT().Get(gomock.Any(), "/activities?urgency_id=456").Return(mockResp, nil)

		result, err := client.GetActivitiesByUrgency(context.Background(), 456)
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "failed to decode activities response")
	})

	t.Run("handles empty activities list", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		listResponse := activityV1.ActivityListResponse{
			Activities: []activityV1.ActivityResponse{},
			Total:      0,
			Page:       1,
			PageSize:   10,
			TotalPages: 0,
		}

		responseBody, _ := json.Marshal(listResponse)
		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}

		mockHTTP.EXPECT().Get(gomock.Any(), "/activities?urgency_id=999").Return(mockResp, nil)

		result, err := client.GetActivitiesByUrgency(context.Background(), 999)
		assert.NoError(t, err)
		require.NotNil(t, result)
		assert.Empty(t, result)
	})
}

func TestActivityClient_LogActivity(t *testing.T) {
	t.Parallel()

	t.Run("successfully logs activity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		expectedResponse := &activityV1.ActivityResponse{
			ID:          1,
			Description: "User logged in",
			EmployeeID:  123,
			UrgencyID:   456,
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		}

		responseBody, _ := json.Marshal(expectedResponse)
		mockResp := &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}

		mockHTTP.EXPECT().Post(gomock.Any(), "/activities", gomock.Any()).DoAndReturn(
			func(ctx context.Context, endpoint string, req *activityV1.ActivityCreateRequest) (*http.Response, error) {
				assert.Equal(t, "User logged in", req.Description)
				assert.Equal(t, uint(123), req.EmployeeID)
				assert.Equal(t, uint(456), req.UrgencyID)
				return mockResp, nil
			})

		err := client.LogActivity(context.Background(), "User logged in", 123, 456)
		assert.NoError(t, err)
	})

	t.Run("returns error when CreateActivity fails", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		mockHTTP.EXPECT().Post(gomock.Any(), "/activities", gomock.Any()).Return(nil, fmt.Errorf("network error"))

		err := client.LogActivity(context.Background(), "Test activity", 123, 456)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create activity")
	})

	t.Run("handles empty description", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		mockResp := &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(bytes.NewReader([]byte(`{"error":"validation failed"}`))),
		}

		mockHTTP.EXPECT().Post(gomock.Any(), "/activities", gomock.Any()).Return(mockResp, nil)

		err := client.LogActivity(context.Background(), "", 123, 456)
		assert.Error(t, err)
	})

	t.Run("handles zero employee and urgency IDs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		expectedResponse := &activityV1.ActivityResponse{
			ID:          1,
			Description: "System activity",
			EmployeeID:  0,
			UrgencyID:   0,
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		}

		responseBody, _ := json.Marshal(expectedResponse)
		mockResp := &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}

		mockHTTP.EXPECT().Post(gomock.Any(), "/activities", gomock.Any()).DoAndReturn(
			func(ctx context.Context, endpoint string, req *activityV1.ActivityCreateRequest) (*http.Response, error) {
				assert.Equal(t, "System activity", req.Description)
				assert.Equal(t, uint(0), req.EmployeeID)
				assert.Equal(t, uint(0), req.UrgencyID)
				return mockResp, nil
			})

		err := client.LogActivity(context.Background(), "System activity", 0, 0)
		assert.NoError(t, err)
	})
}

func TestNewActivityClient(t *testing.T) {
	t.Parallel()

	t.Run("creates client with valid configuration", func(t *testing.T) {
		logger := utils.NewTestLogger()
		serviceAuth := auth.NewServiceAuth(auth.ServiceAuthConfig{
			Secret:      "test-secret",
			ServiceName: "urgency-service",
			TokenTTL:    time.Hour,
		})
		config := ActivityClientConfig{
			BaseURL:     "http://localhost:8080",
			ServiceAuth: serviceAuth,
			Logger:      logger,
			Timeout:     30 * time.Second,
		}

		client := NewActivityClient(config)
		assert.NotNil(t, client)

		// Verify it implements the interface
		var _ ActivityClient = client
	})

	t.Run("creates client with minimal configuration", func(t *testing.T) {
		logger := utils.NewTestLogger()
		config := ActivityClientConfig{
			BaseURL: "http://activity-service:8080",
			Logger:  logger,
		}

		client := NewActivityClient(config)
		assert.NotNil(t, client)
	})

	t.Run("creates client with custom timeout", func(t *testing.T) {
		logger := utils.NewTestLogger()
		config := ActivityClientConfig{
			BaseURL: "http://activity-service:8080",
			Logger:  logger,
			Timeout: 60 * time.Second,
		}

		client := NewActivityClient(config)
		assert.NotNil(t, client)
	})
}

func TestActivityClient_EdgeCases(t *testing.T) {
	t.Parallel()

	t.Run("handles context cancellation", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		req := &activityV1.ActivityCreateRequest{
			Description: "Test activity",
			EmployeeID:  123,
			UrgencyID:   456,
		}

		mockHTTP.EXPECT().Post(ctx, "/activities", req).Return(nil, context.Canceled)

		result, err := client.CreateActivity(ctx, req)
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("handles large urgency ID", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		largeID := uint(999999999)
		listResponse := activityV1.ActivityListResponse{
			Activities: []activityV1.ActivityResponse{},
			Total:      0,
			Page:       1,
			PageSize:   10,
			TotalPages: 0,
		}

		responseBody, _ := json.Marshal(listResponse)
		mockResp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}

		expectedEndpoint := fmt.Sprintf("/activities?urgency_id=%d", largeID)
		mockHTTP.EXPECT().Get(gomock.Any(), expectedEndpoint).Return(mockResp, nil)

		result, err := client.GetActivitiesByUrgency(context.Background(), largeID)
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("handles very long description", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockHTTP := NewMockHTTPClient(ctrl)
		logger := utils.NewTestLogger()

		client := &activityClient{
			httpClient: mockHTTP,
			logger:     logger.WithName("activityClient"),
		}

		longDescription := string(make([]byte, 1000))
		for i := range longDescription {
			longDescription = longDescription[:i] + "a" + longDescription[i+1:]
		}

		expectedResponse := &activityV1.ActivityResponse{
			ID:          1,
			Description: longDescription,
			EmployeeID:  123,
			UrgencyID:   456,
			CreatedAt:   time.Now().Format(time.RFC3339),
			UpdatedAt:   time.Now().Format(time.RFC3339),
		}

		responseBody, _ := json.Marshal(expectedResponse)
		mockResp := &http.Response{
			StatusCode: http.StatusCreated,
			Body:       io.NopCloser(bytes.NewReader(responseBody)),
		}

		mockHTTP.EXPECT().Post(gomock.Any(), "/activities", gomock.Any()).DoAndReturn(
			func(ctx context.Context, endpoint string, req *activityV1.ActivityCreateRequest) (*http.Response, error) {
				assert.Equal(t, longDescription, req.Description)
				return mockResp, nil
			})

		err := client.LogActivity(context.Background(), longDescription, 123, 456)
		assert.NoError(t, err)
	})
}
