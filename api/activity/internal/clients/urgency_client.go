package clients

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
    "github.com/pd120424d/mountain-service/api/shared/auth"
    "github.com/pd120424d/mountain-service/api/shared/client"
    "github.com/pd120424d/mountain-service/api/shared/utils"
)

// UrgencyClient provides access to the urgency service for internal S2S calls
//go:generate mockgen -source=urgency_client.go -destination=urgency_client_gomock.go -package=clients mountain_service/activity/internal/clients -imports=gomock=go.uber.org/mock/gomock -typed

type UrgencyClient interface {
    GetUrgencyByID(ctx context.Context, id uint) (*urgencyV1.UrgencyResponse, error)
}

type HTTPClient interface {
    Get(ctx context.Context, endpoint string) (*http.Response, error)
}

type urgencyClient struct {
    httpClient HTTPClient
    logger     utils.Logger
}

type UrgencyClientConfig struct {
    BaseURL     string
    ServiceAuth auth.ServiceAuth
    Logger      utils.Logger
    Timeout     time.Duration
}

func NewUrgencyClient(config UrgencyClientConfig) UrgencyClient {
    httpClient := client.NewHTTPClient(client.HTTPClientConfig{
        BaseURL:     config.BaseURL,
        ServiceAuth: config.ServiceAuth,
        Logger:      config.Logger,
        Timeout:     config.Timeout,
    })

    return &urgencyClient{httpClient: httpClient, logger: config.Logger.WithName("urgencyClient")}
}

func (c *urgencyClient) GetUrgencyByID(ctx context.Context, id uint) (*urgencyV1.UrgencyResponse, error) {
    endpoint := fmt.Sprintf("/api/v1/service/urgency/%d", id)
    resp, err := c.httpClient.Get(ctx, endpoint)
    if err != nil {
        return nil, fmt.Errorf("failed to call urgency service: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("urgency service returned status %d", resp.StatusCode)
    }

    var ur urgencyV1.UrgencyResponse
    if err := json.NewDecoder(resp.Body).Decode(&ur); err != nil {
        return nil, fmt.Errorf("failed to decode urgency response: %w", err)
    }
    return &ur, nil
}

