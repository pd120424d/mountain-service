package urgency

import (
	"context"
	"encoding/json"
	"fmt"
	"net"

	"net/http"
	"os"
	"time"

	urgencyV1 "github.com/pd120424d/mountain-service/api/contracts/urgency/v1"
	"github.com/pd120424d/mountain-service/api/shared/auth"
	"github.com/pd120424d/mountain-service/api/shared/client"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

// Client defines the S2S client for the urgency service.
type Client interface {
	GetUrgencyByID(ctx context.Context, id uint) (*urgencyV1.UrgencyResponse, error)
}

// Config for constructing a Client.
type Config struct {
	BaseURL     string
	ServiceAuth auth.ServiceAuth
	Logger      utils.Logger
	Timeout     time.Duration
}

type httpClient interface {
	Get(ctx context.Context, endpoint string) (*http.Response, error)
}

type clientImpl struct {
	http   httpClient
	logger utils.Logger
	// simple retries for transient errors
	maxRetries int
}

// New creates a new urgency S2S client using the shared HTTP client.
func New(cfg Config) Client {
	if cfg.Timeout == 0 {
		cfg.Timeout = 30 * time.Second
	}
	c := client.NewHTTPClient(client.HTTPClientConfig{
		BaseURL:     cfg.BaseURL,
		Timeout:     cfg.Timeout,
		ServiceAuth: cfg.ServiceAuth,
		Logger:      cfg.Logger,
	})
	return &clientImpl{http: c, logger: cfg.Logger.WithName("urgencyS2S"), maxRetries: 2}
}

// NewFromEnv constructs a Client using URGENCY_SERVICE_URL or a sane in-cluster default.
func NewFromEnv(logger utils.Logger, sa auth.ServiceAuth) Client {
	base := os.Getenv("URGENCY_SERVICE_URL")
	if base == "" {
		base = "http://urgency-service:8083"
	}
	return New(Config{BaseURL: base, ServiceAuth: sa, Logger: logger, Timeout: 30 * time.Second})
}

func (c *clientImpl) GetUrgencyByID(ctx context.Context, id uint) (*urgencyV1.UrgencyResponse, error) {
	log := c.logger.WithContext(ctx)
	endpoint := fmt.Sprintf("/api/v1/service/urgency/%d", id)
	resp, err := c.retryGet(ctx, endpoint)
	if err != nil {
		log.Errorf("urgency.get_by_id http_error id=%d err=%v", id, err)
		return nil, fmt.Errorf("failed to call urgency service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		log.Warnf("urgency.get_by_id not_found id=%d", id)
		return nil, fmt.Errorf("urgency %d not found", id)
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf("urgency.get_by_id non_200 id=%d status=%d", id, resp.StatusCode)
		return nil, fmt.Errorf("urgency service returned status %d", resp.StatusCode)
	}

	var ur urgencyV1.UrgencyResponse
	if err := json.NewDecoder(resp.Body).Decode(&ur); err != nil {
		log.Errorf("urgency.get_by_id decode_error id=%d err=%v", id, err)
		return nil, fmt.Errorf("failed to decode urgency response: %w", err)
	}
	log.Infof("urgency.get_by_id success id=%d status=%s", id, ur.Status)
	return &ur, nil
}

func (c *clientImpl) retryGet(ctx context.Context, endpoint string) (*http.Response, error) {
	var lastErr error
	backoff := 100 * time.Millisecond
	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		resp, err := c.http.Get(ctx, endpoint)
		if err == nil {
			if resp.StatusCode == http.StatusTooManyRequests || (resp.StatusCode >= 500 && resp.StatusCode <= 599) {
				lastErr = fmt.Errorf("server status %d", resp.StatusCode)
				resp.Body.Close()
			} else {
				return resp, nil
			}
		} else {
			if ne, ok := err.(net.Error); ok && (ne.Timeout() || ne.Temporary()) {
				lastErr = err
			} else if err == context.DeadlineExceeded || err == context.Canceled {
				return nil, err
			} else {
				lastErr = err
			}
		}
		if attempt == c.maxRetries {
			break
		}
		select {
		case <-time.After(backoff):
			backoff = time.Duration(float64(backoff) * 1.8)
			continue
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	return nil, lastErr
}
