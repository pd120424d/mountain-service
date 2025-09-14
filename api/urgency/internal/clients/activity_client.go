package clients

//go:generate mockgen -source=activity_client.go -destination=activity_client_gomock.go -package=clients mountain_service/urgency/internal/clients -imports=gomock=go.uber.org/mock/gomock -typed

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	activityV1 "github.com/pd120424d/mountain-service/api/contracts/activity/v1"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/client"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

// HTTPClient abstracts the shared HTTP client for easier testing.
type HTTPClient interface {
	Get(ctx context.Context, endpoint string) (*http.Response, error)
	Post(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)
	Put(ctx context.Context, endpoint string, body interface{}) (*http.Response, error)
	Delete(ctx context.Context, endpoint string) (*http.Response, error)
}

// ActivityClient interface for communicating with the activity service
type ActivityClient interface {
	CreateActivity(ctx context.Context, req *activityV1.ActivityCreateRequest) (*activityV1.ActivityResponse, error)
	GetActivitiesByUrgency(ctx context.Context, urgencyID uint) ([]activityV1.ActivityResponse, error)
	LogActivity(ctx context.Context, description string, employeeID, urgencyID uint) error
}

// activityClient implements the ActivityClient interface
type activityClient struct {
	httpClient HTTPClient
	logger     utils.Logger
}

// ActivityClientConfig holds the configuration for the activity client
type ActivityClientConfig struct {
	BaseURL     string
	ServiceAuth auth.ServiceAuth
	Logger      utils.Logger
	Timeout     time.Duration
}

// NewActivityClient creates a new activity service client
func NewActivityClient(config ActivityClientConfig) ActivityClient {
	httpClient := client.NewHTTPClient(client.HTTPClientConfig{
		BaseURL:     config.BaseURL,
		ServiceAuth: config.ServiceAuth,
		Logger:      config.Logger,
		Timeout:     config.Timeout,
	})

	return &activityClient{
		httpClient: httpClient,
		logger:     config.Logger.WithName("activityClient"),
	}
}

// CreateActivity creates a new activity
func (c *activityClient) CreateActivity(ctx context.Context, req *activityV1.ActivityCreateRequest) (*activityV1.ActivityResponse, error) {
	log := c.logger.WithContext(ctx)
	log.Infof("Creating activity for urgency %d by employee %d", req.UrgencyID, req.EmployeeID)

	resp, err := c.httpClient.Post(ctx, "/api/v1/service/activities", req)
	if err != nil {
		log.Errorf("Failed to create activity: %v", err)
		return nil, fmt.Errorf("failed to create activity: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		log.Errorf("Activity service returned status %d", resp.StatusCode)
		return nil, fmt.Errorf("activity service returned status %d", resp.StatusCode)
	}

	var activity activityV1.ActivityResponse
	if err := json.NewDecoder(resp.Body).Decode(&activity); err != nil {
		log.Errorf("Failed to decode activity response: %v", err)
		return nil, fmt.Errorf("failed to decode activity response: %w", err)
	}

	log.Infof("Successfully created activity with ID %d", activity.ID)
	return &activity, nil
}

// GetActivitiesByUrgency retrieves all activities for a specific urgency
func (c *activityClient) GetActivitiesByUrgency(ctx context.Context, urgencyID uint) ([]activityV1.ActivityResponse, error) {
	log := c.logger.WithContext(ctx)
	log.Infof("Getting activities for urgency %d", urgencyID)

	endpoint := fmt.Sprintf("/api/v1/service/activities?urgencyId=%d", urgencyID)
	resp, err := c.httpClient.Get(ctx, endpoint)
	if err != nil {
		log.Errorf("Failed to get activities: %v", err)
		return nil, fmt.Errorf("failed to get activities: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Errorf("Activity service returned status %d", resp.StatusCode)
		return nil, fmt.Errorf("activity service returned status %d", resp.StatusCode)
	}

	var listResponse activityV1.ActivityListResponse
	if err := json.NewDecoder(resp.Body).Decode(&listResponse); err != nil {
		log.Errorf("Failed to decode activities response: %v", err)
		return nil, fmt.Errorf("failed to decode activities response: %w", err)
	}

	log.Infof("Successfully retrieved %d activities for urgency %d", len(listResponse.Activities), urgencyID)
	return listResponse.Activities, nil
}

// LogActivity is a convenience method for logging activities
func (c *activityClient) LogActivity(ctx context.Context, description string, employeeID, urgencyID uint) error {
	req := &activityV1.ActivityCreateRequest{
		Description: description,
		EmployeeID:  employeeID,
		UrgencyID:   urgencyID,
	}

	_, err := c.CreateActivity(ctx, req)
	return err
}
